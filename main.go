package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"golang.org/x/sys/unix"
)

func client(cid, port uint32, c cred) {
	fd, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}

	sa := &unix.SockaddrVM{
		CID:  cid,
		Port: uint32(port),
	}

	if err := unix.Connect(fd, sa); err != nil {
		panic(err)
	}

	m := msg{
		AccessKeyID:     c.AccessKeyID,
		SecretAccessKey: c.SecretAccessKey,
		Token:           c.Token,
		Ciphertext:      "xddgang",
	}

	b, _ := json.Marshal(m)

	if err := unix.Send(fd, b, 0); err != nil {
		panic(err)
	}
}

const (
	backlog = 128
	bufsize = 4096
)

type cred struct {
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	Token           string `json:"Token"`
}

func server(port uint32) {
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

			var m msg
			if err := json.Unmarshal(buf[:n], &m); err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				break
			}

			cmd := exec.Command(
				"./kmstool_enclave_cli",
				"decrypt",
				"--region", "us-east-2",
				"--proxy-port", "8000",
				"--aws-access-key-id", m.AccessKeyID,
				"--aws-secret-access-key", m.SecretAccessKey,
				"--aws-session-token", m.Token,
				"--ciphertext", m.Ciphertext,
			)

			out, err := cmd.Output()
			if err != nil {
				log.Printf("Failed executing KMS command: %s", err)
				if err, ok := err.(*exec.ExitError); ok {
					log.Printf("Stderr: %s", string(err.Stderr))
				}
			} else {
				log.Printf("Got message: %s.", string(out))
			}
		}
	}
}

type msg struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Token           string `json:"token"`
	Ciphertext      string `json:"ciphertext"`
}

func main() {
	if len(os.Args) == 1 {
		panic("subcommand client or server required")
	}

	switch os.Args[1] {
	case "client":
		ctx := context.TODO()

		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		md := imds.NewFromConfig(cfg)
		mo, err := md.GetMetadata(ctx, &imds.GetMetadataInput{Path: "iam/security-credentials/dev-ec2-test-enclave"})
		if err != nil {
			panic(err)
		}

		defer mo.Content.Close()

		b, err := io.ReadAll(mo.Content)
		if err != nil {
			panic(err)
		}

		var c cred
		if err := json.Unmarshal(b, &c); err != nil {
			panic(err)
		}

		a := os.Args[2:]
		if len(a) != 2 {
			panic("cid and port arguments required")
		}
		cid, err := strconv.ParseUint(a[0], 10, 32)
		if err != nil {
			panic(err)
		}
		port, err := strconv.ParseUint(a[1], 10, 32)
		if err != nil {
			panic(err)
		}
		client(uint32(cid), uint32(port), c)
	case "server":
		a := os.Args[2:]
		if len(a) != 1 {
			panic("port argument required")
		}
		port, err := strconv.ParseUint(a[0], 10, 32)
		if err != nil {
			panic(err)
		}

		log.Printf("Starting server on port %d.", port)
		server(uint32(port))
	default:
		panic("unrecognized subcommand")
	}
}
