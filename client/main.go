package main

import (
	"golang.org/x/sys/unix"
)

const clientCID = 3
const backlog = 128
const port = 5005
const bufsize = 1024

func main() {
	fd, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}

	sa := &unix.SockaddrVM{
		CID:  unix.VMADDR_CID_ANY,
		Port: port,
	}

	if err := unix.Connect(fd, sa); err != nil {
		panic(err)
	}

	if err := unix.Send(fd, []byte("Hello, sir."), 0); err != nil {
		panic(err)
	}
}
