package main

import (
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

func client(cid, port uint32) {
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

	if err := unix.Send(fd, []byte("Hello, world!"), 0); err != nil {
		panic(err)
	}
}

func server(port uint32) {
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

	if err := unix.Send(fd, []byte("Hello, world!"), 0); err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) == 1 {
		panic("subcommand client or server required")
	}

	switch os.Args[1] {
	case "client":
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
		client(uint32(cid), uint32(port))
	case "server":
		a := os.Args[2:]
		if len(a) != 1 {
			panic("port argument required")
		}
		port, err := strconv.ParseUint(a[0], 10, 32)
		if err != nil {
			panic(err)
		}
		server(uint32(port))
	default:
		panic("unrecognized subcommand")
	}
}
