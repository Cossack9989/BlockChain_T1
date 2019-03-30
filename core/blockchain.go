package core

import (
	"fmt"
	"github.com/boltdb/bolt"
)

const dbFile = "Blocks.db"
const blocksBucket = "Blocks"

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func NewBlockChain() *BlockChain {
	var tip []byte
	db, _ := bolt.Open(dbFile, 0600, nil)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
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

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	newBlock := NewBlock(data, lastHash)
	bc.db.Update(func(tx *bolt.Tx) error {
		hash := newBlock.GetBlockHash()
		b := tx.Bucket([]byte(blocksBucket))
		b.Put(hash, newBlock.Serialize())
		b.Put([]byte("l"), hash)
		if 1 == 1 {
			fmt.Printf("Prev's hash: %x\n", newBlock.PrevBlockHash)
			fmt.Printf("    Data   : %s\n", newBlock.Data)
			fmt.Printf("    Proof  : %d\n", newBlock.Proof)
			fmt.Println()
		}
		return nil
	})

}
