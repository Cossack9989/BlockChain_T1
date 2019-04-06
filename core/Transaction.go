package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
)

const subsidy = 50

type Transaction struct {
	ID   []byte     `json:"id"`
	Vin  []TXInput  `json:"Vin"`
	Vout []TXOutput `json:"Vout"`
}

func NewCoinbaseTX(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(subsidy, to)
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
	ws := bc.Wallets
	w := ws.GetWallet(from)
	var inputs []TXInput
	var outputs []TXOutput
	if acc < amount {
		fmt.Println("GG")
	}
	for id, outs := range validOutputs {
		txid := []byte(id)
		for _, out := range outs {
			inputs = append(inputs, TXInput{txid, out, nil, w.PublicKey})
		}
	}
	outputs = append(outputs, NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, NewTXOutput(acc-amount, from))
	}
	tx := &Transaction{nil, inputs, outputs}
	tx.ID = tx.GetHash()
	bc.SignTransaction(tx, w.PrivateKey)
	return tx
}

func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, previousTXs map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}
	txCopy := tx.TrimmedCopy()
	for inID, vin := range txCopy.Vin {
		prevTx := previousTXs[string(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.GetHash()
		txCopy.Vin[inID].PubKey = nil
		r, s, _ := ecdsa.Sign(rand.Reader, &privateKey, txCopy.ID)
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}
	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[string(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.GetHash()
		txCopy.Vin[inID].PubKey = nil
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}
