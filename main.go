package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/DIMO-Network/synthetic-wallet-enclave/pkg/types"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"golang.org/x/sys/unix"
)

const (
	backlog = 128
	bufsize = 4096
)

var seed []byte
var seedMu sync.RWMutex

func handle(buf []byte, logger *zerolog.Logger) (res []byte, err error) {
	var m types.Request[json.RawMessage]
	if err := json.Unmarshal(buf, &m); err != nil {
		return nil, err
	}

	seedMu.RLock()

	for seed == nil {
		logger.Info().Msg("Seed not populated, plan to call KMS.")
		seedMu.RUnlock()
		if err := func() error {
			seedMu.Lock()
			defer seedMu.Unlock()
			if seed == nil {
				cmd := exec.Command(
					"./kmstool_enclave_cli",
					"decrypt",
					"--region", "us-east-2",
					"--proxy-port", "8000",
					"--aws-access-key-id", m.Credentials.AccessKeyID,
					"--aws-secret-access-key", m.Credentials.SecretAccessKey,
					"--aws-session-token", m.Credentials.Token,
					"--ciphertext", m.EncryptedSeed,
				)

				out, err := cmd.Output()
				if err != nil {
					return err
				}

				// Output has the form
				// PLAINTEXT: <base64-encoded plaintext>
				seed, err = base64.StdEncoding.DecodeString(strings.TrimSpace(strings.Split(string(out), ":")[1]))
				if err != nil {
					return err
				}
			}
			return nil
		}(); err != nil {
			return nil, err
		}
		seedMu.RLock()
	}

	ek, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	seedMu.RUnlock()
	if err != nil {
		return nil, err
	}

	switch m.Type {
	case "GetAddress":
		var x types.AddrReqData
		if err := json.Unmarshal(m.Data, &x); err != nil {
			return nil, err
		}

		ck, err := ek.Child(hdkeychain.HardenedKeyStart + x.ChildNumber)
		if err != nil {
			return nil, err
		}

		add, err := ck.Address(&chaincfg.MainNetParams)
		if err != nil {
			return nil, err
		}

		addr := common.BytesToAddress(add.ScriptAddress())

		return json.Marshal(types.Response[types.AddrResData]{Code: 0, Data: types.AddrResData{Address: addr}})
	case "SignHash":
		var x types.SignReqData
		if err := json.Unmarshal(m.Data, &x); err != nil {
			return nil, err
		}

		ck, err := ek.Child(hdkeychain.HardenedKeyStart + x.ChildNumber)
		if err != nil {
			return nil, err
		}

		pk, err := ck.ECPrivKey()
		if err != nil {
			return nil, err
		}

		sig, err := crypto.Sign(x.Hash.Bytes(), pk.ToECDSA())
		if err != nil {
			return nil, err
		}

		sig[64] += 27

		return json.Marshal(types.Response[types.SignResData]{Code: 0, Data: types.SignResData{Signature: sig}})
	default:
		return nil, fmt.Errorf("unrecognized request type %s", m.Type)
	}
}

func accept(fd int, logger *zerolog.Logger) error {
	nfd, _, err := unix.Accept(fd)
	if err != nil {
		return err
	}
	defer unix.Close(nfd)

	buf := make([]byte, bufsize)
	n, _, err := unix.Recvfrom(nfd, buf, 0)
	if err != nil {
		return err
	}

	logger.Debug().Msgf("Got message %q.", string(buf[:n]))

	res, err := handle(buf[:n], logger)
	if err != nil {
		res, _ = json.Marshal(types.Response[types.ErrData]{Code: 2, Data: types.ErrData{Message: err.Error()}})
	}

	return unix.Send(nfd, res, 0)
}

func enclave(ctx context.Context, port uint32, logger *zerolog.Logger) error {
	fd, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
	if err != nil {
		return err
	}

	logger.Debug().Msgf("Created socket %d.", fd)

	sa := &unix.SockaddrVM{
		CID:  unix.VMADDR_CID_ANY,
		Port: port,
	}

	if err := unix.Bind(fd, sa); err != nil {
		panic(err)
	}

	logger.Debug().Msgf("Bound socket with a random address and port %d.", port)

	if err := unix.Listen(fd, backlog); err != nil {
		panic(err)
	}

	logger.Debug().Msgf("Accepting requests with backlog %d.", backlog)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := accept(fd, logger); err != nil {
				logger.Err(err).Msg("Accept failed.")
			}
		}
	}
}

func main() {
	logger := zerolog.New(os.Stderr).With().Str("app", "virtual-device-enclave").Timestamp().Logger()

	if len(os.Args) < 2 {
		logger.Fatal().Msg("Port argument required.")
	}
	port, err := parseUint32(os.Args[1])
	if err != nil {
		logger.Fatal().Err(err).Msgf("Couldn't parse port %q.", os.Args[1])
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	enclave(ctx, port, &logger)
}

func parseUint32(s string) (uint32, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(n), err
}
