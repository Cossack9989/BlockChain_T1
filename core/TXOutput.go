package core

import ()

type TXOutput struct {
	Value        int    `json:"value"`
	ScriptPubKey string `json:"scriptPubKey"`
}

func (out TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
