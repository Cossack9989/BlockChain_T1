package core

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

const BlocksFile = "Blocks.db"
const blocksBucket = "Blocks"
const dbgg = 0

type BlockChain struct {
	tip     []byte
	db      *bolt.DB
	Wallets *Wallets
}

func NewBlockChain(pubKeyHash string) *BlockChain {
	var tip []byte
	db, _ := bolt.Open(BlocksFile, 0600, nil)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			cb := NewCoinbaseTX(pubKeyHash, "The Beginer!")
			genesis := NewGenesisBlock(cb)
			b, _ := tx.CreateBucket([]byte(blocksBucket))
			hash := genesis.GetBlockHash()
			b.Put(hash, genesis.Serialize())
			b.Put([]byte("l"), hash)

			tip = hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	ws := NewWallets()
	bc := &BlockChain{tip, db, ws}
	return bc
}

func (bc *BlockChain) AddBlock(txs []*Transaction) {
	lastHash := bc.tip
	newBlock := NewBlock(txs, lastHash)
	bc.db.Update(func(tx *bolt.Tx) error {
		hash := newBlock.GetBlockHash()
		b := tx.Bucket([]byte(blocksBucket))
		b.Put(hash, newBlock.Serialize())
		b.Put([]byte("l"), hash)
		bc.tip = hash
		if dbgg == 1 {
			fmt.Printf("Prev's hash: %x\n", newBlock.PrevBlockHash)
			fmt.Printf("    Proof  : %d\n", newBlock.Proof)
			fmt.Println()
		}
		return nil
	})
}

func (bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []*Transaction {
	var unspentTXs []*Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.GetIterator()
	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txid := string(tx.ID)
		OutPuts:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txid] != nil {
					for _, spentOut := range spentTXOs[txid] {
						if spentOut == outIdx {
							continue OutPuts
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, tx)
				}
			}
			if tx.IsCoinBase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockWithPubkey(pubKeyHash) {
						inTxId := string(in.Txid)
						spentTXOs[inTxId] = append(spentTXOs[inTxId], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

func (bc *BlockChain) FindUnspentTransactionsOuts(pubKeyHash []byte) []TXOutput {
	var unspentTXsOuts []TXOutput
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	for _, txs := range unspentTXs {
		for _, out := range txs.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				unspentTXsOuts = append(unspentTXsOuts, out)
			}
		}
	}
	return unspentTXsOuts
}

func (bc *BlockChain) FindSpendableOutputs(from string, amount int) (int, map[string][]int) {
	pubkeyhash := Base58Decode([]byte(from))
	pubkeyhash = pubkeyhash[0 : len(pubkeyhash)-4]
	unspentTXs := bc.FindUnspentTransactions(pubkeyhash)
	unspentTXOuts := make(map[string][]int)
	accumulated := 0
find:
	for _, txs := range unspentTXs {
		idtx := string(txs.ID)
		for ido, out := range txs.Vout {
			if out.IsLockedWithKey(pubkeyhash) {
				accumulated += out.Value
				unspentTXOuts[idtx] = append(unspentTXOuts[idtx], ido)
				if accumulated >= amount {
					break find
				}
			}
		}
	}
	return accumulated, unspentTXOuts
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	previousTXs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		previousTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			fmt.Println(err)
		} else {
			previousTXs[string(previousTX.ID)] = previousTX
		}
	}
	tx.Sign(privateKey, previousTXs)
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.GetIterator()
	for {
		b := bci.Next()
		for _, tx := range b.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(b.PrevBlockHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("Transaction is not found")
}
