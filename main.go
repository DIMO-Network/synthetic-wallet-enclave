package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sys/unix"
)

const (
	backlog = 128
	bufsize = 4096
)

func enclave(port uint32) {
	fd, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}

	log.Printf("Opened file descriptor %d.", fd)

	sa := &unix.SockaddrVM{
		CID:  unix.VMADDR_CID_ANY,
		Port: port,
	}

	if err := unix.Bind(fd, sa); err != nil {
		panic(err)
	}

	if err := unix.Listen(fd, backlog); err != nil {
		panic(err)
	}

	buf := make([]byte, bufsize)

	for {
		nfd, _, err := unix.Accept(fd)
		if err != nil {
			log.Printf("Error on accept: %s.", err)
			continue
		}

		for {
			n, _, err := unix.Recvfrom(nfd, buf, 0)
			if err != nil {
				log.Printf("Error on recvfrom: %s.", err)
				break
			}
			if n == 0 {
				break
			}

			log.Printf("Got message: %s", string(buf[:n]))

			var m Request
			if err := json.Unmarshal(buf[:n], &m); err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				break
			}

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
				log.Printf("Failed executing KMS command: %s", err)
				if err, ok := err.(*exec.ExitError); ok {
					log.Printf("Stderr: %s", string(err.Stderr))
				}
			} else {
				sout := string(out)
				plain := strings.Split(sout, ":")[1]

				seed, err := base64.StdEncoding.DecodeString(plain)
				if err != nil {
					log.Println(err)
					continue
				}

				ek, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
				if err != nil {
					log.Println(err)
					continue
				}

				ck, err := ek.Child(hdkeychain.HardenedKeyStart + m.ChildNumber)
				if err != nil {
					log.Println(err)
					continue
				}

				add, err := ck.Address(&chaincfg.MainNetParams)
				if err != nil {
					log.Println(err)
					continue
				}

				addr := common.BytesToAddress(add.ScriptAddress())

				bout, err := json.Marshal(Response{Address: addr})
				if err != nil {
					log.Println(err)
					continue
				}

				err = unix.Send(nfd, bout, 0)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

type Request struct {
	Credentials   AWSCredentials `json:"credentials"`
	EncryptedSeed string         `json:"encryptedSeed"`
	ChildNumber   uint32         `json:"childNumber"`
}

type Response struct {
	Address common.Address `json:"address"`
}

type AWSCredentials struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Token           string `json:"token"`
}

func main() {
	if len(os.Args) < 2 {
		panic("port argument required")
	}
	port, err := parseUint32(os.Args[1])
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on port %d.", port)
	enclave(port)
}

func parseUint32(s string) (uint32, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(n), err
}
