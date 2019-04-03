package core

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

const subsidy = 1

type Transaction struct {
	ID   []byte     `json:"id"`
	Vin  []TXInput  `json:"Vin"`
	Vout []TXOutput `json:"Vout"`
}

func NewCoinbaseTX(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := &Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.ID = tx.GetHash()
	return tx
}

func (tx *Transaction) GetHash() []byte {
	txCopy := &Transaction{[]byte{}, tx.Vin, tx.Vout}
	js := txCopy.Serialize()
	hash := sha256.Sum256(js)
	return hash[:]
}

func (tx *Transaction) Serialize() []byte {
	result, _ := json.Marshal(tx)
	return []byte(result)
}

func DeserializeTransaction(data []byte) *Transaction {
	var tx *Transaction
	json.Unmarshal(data, tx)
	return tx
}
func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (bc *BlockChain) NewTransaction(from string, to string, amount int) *Transaction {
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	var inputs []TXInput
	var outputs []TXOutput
	if acc < amount {
		fmt.Println("GG")
	}
	for id, outs := range validOutputs {
		txid := []byte(id)
		for _, out := range outs {
			inputs = append(inputs, TXInput{txid, out, from})
		}
	}
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}
	tx := &Transaction{nil, inputs, outputs}
	tx.ID = tx.GetHash()
	return tx
}
