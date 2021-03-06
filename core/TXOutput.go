package core

import (
	"bytes"
)

type TXOutput struct {
	Value      int    `json:"value"`
	PubKeyHash []byte `json:"pubKeyHash"`
}

func NewTXOutput(value int, address string) TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return *txo
}

func (out *TXOutput) Lock(address []byte) {
	pubkeyhash := Base58Decode(address)
	pubkeyhash = pubkeyhash[0 : len(pubkeyhash)-4]
	out.PubKeyHash = pubkeyhash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(pubKeyHash, out.PubKeyHash) == 0
}
