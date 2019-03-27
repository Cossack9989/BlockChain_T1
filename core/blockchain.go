package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deckarep/golang-set"
)

type Block struct {
	Index        int           `json:"index"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
}
type Transaction struct {
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
}
type Chain struct {
	Chain  []Block `json:"chain"`
	Length int     `json:"length"`
}
type Blockchain struct {
	CurrentTransactions []Transaction
	Chain               []Block
	Nodes               mapset.Set
}

func (t *Blockchain) Genesis() *Blockchain {
	//balabala
	return t
}
func (t *Blockchain) NewBlock(proof int, previousHash string) Block {
	block := new(Block)
	//balabala
	return *block
}
func (t *Blockchain) NewTransaction(sender string, recipient string, amount int) int {
	transaction := new(Transaction)
	//balabala
	return 0
}
