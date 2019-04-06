package core

import (
	"bytes"
)

type TXInput struct {
	Txid []byte `json:"Txid"`
	Vout int    `json:"Vout"`
	Signature []byte `json:"Signature"`
	PubKey    []byte `json:"pubKey"`
}

func (in *TXInput) CanUnlockWithPubkey(pubKeyHash []byte) bool {
	lockingKey := HashPubKey(in.PubKey)
	return bytes.Compare(lockingKey, pubKeyHash) == 0
}
