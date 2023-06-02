package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Request[A any] struct {
	Credentials struct {
		AccessKeyID     string `json:"accessKeyId"`
		SecretAccessKey string `json:"secretAccessKey"`
		Token           string `json:"token"`
	} `json:"credentials"`
	Type          string `json:"type"`
	EncryptedSeed string `json:"encryptedSeed"`
	Data          A      `json:"data"`
}

type AddrReqData struct {
	ChildNumber uint32 `json:"childNumber"`
}

type SignReqData struct {
	ChildNumber uint32      `json:"childNumber"`
	Hash        common.Hash `json:"hash"`
}

type AddrResData struct {
	Address common.Address `json:"address"`
}

type SignResData struct {
	Signature hexutil.Bytes `json:"signature"`
}

type ErrData struct {
	Message string `json:"message"`
}

type Response[A any] struct {
	Code int `json:"code"`
	Data A   `json:"data"`
}
