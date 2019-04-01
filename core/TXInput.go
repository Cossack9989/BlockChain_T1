package core

import ()

type TXInput struct {
	Txid      []byte `json:"Txid"`
	Vout      int    `json:"Vout"`
	ScriptSig string `json:"scriptSig"`
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}
