package core

import ()

type TXInput struct {
	Txid      []byte `json:"Txid"`
	Vout      int    `json:"Vout"`
	ScriptSig string `json:"sctiptSig"`
}
