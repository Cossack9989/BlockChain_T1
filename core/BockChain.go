package core

import (
	"fmt"
	"github.com/boltdb/bolt"
)

const dbFile = "Blocks.db"
const blocksBucket = "Blocks"
const dbgg = 0

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func NewBlockChain(address string) *BlockChain {
	var tip []byte
	db, _ := bolt.Open(dbFile, 0600, nil)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			cb := NewCoinbaseTX(address, "The Beginer!")
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
	bc := &BlockChain{tip, db}
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

func (bc *BlockChain) FindUnspentTransactions(address string) []*Transaction {
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
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, tx)
				}
			}
			if tx.IsCoinBase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
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

func (bc *BlockChain) FindUnspentTransactionsOuts(address string) []TXOutput {
	var unspentTXsOuts []TXOutput
	unspentTXs := bc.FindUnspentTransactions(address)
	for _, txs := range unspentTXs {
		for _, out := range txs.Vout {
			if out.CanBeUnlockedWith(address) {
				unspentTXsOuts = append(unspentTXsOuts, out)
			}
		}
	}
	return unspentTXsOuts
}

func (bc *BlockChain) FindSpendableOutputs(from string, amount int) (int, map[string][]int) {
	unspentTXs := bc.FindUnspentTransactions(from)
	unspentTXOuts := make(map[string][]int)
	accumulated := 0
find:
	for _, txs := range unspentTXs {
		idtx := string(txs.ID)
		for ido, out := range txs.Vout {
			if out.CanBeUnlockedWith(from) {
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
