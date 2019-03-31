package core

import ()

type TXOutput struct {
	Value        int    `json:"value"`
	ScriptPubKey string `json:"scriptPubKey"`
}
