package core

import (
	"bytes"
)

type TXInput struct {
	Txid []byte `json:"Txid"`
	Vout int    `json:"Vout"`
	// ScriptSig string `json:"scriptSig"`
	Signature []byte `json:"Signature"`
	PubKey    []byte `json:"pubKey"`
}

// func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
// 	return in.ScriptSig == unlockingData
// }

func (in *TXInput) UseKey(pubKeyHash []byte) bool {
	lockingKey := HashPubKey(in.PubKey)
	return bytes.Compare(lockingKey, pubKeyHash) == 0
}
