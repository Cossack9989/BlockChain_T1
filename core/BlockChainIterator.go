package core

import (
	"bytes"
	"github.com/boltdb/bolt"
)

type BlockChainIterator struct {
	CurrentHash []byte
	db          *bolt.DB
}

func (bc *BlockChain) GetIterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}
	return bci
}

func (bci *BlockChainIterator) Next() *Block {
	var currentBlock *Block
	bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		jsstr := b.Get(bci.CurrentHash)
		currentBlock = DeserializeBlock(jsstr)

		return nil
	})
	if !bytes.Equal(currentBlock.PrevBlockHash, []byte{}) {
		bci.CurrentHash = currentBlock.PrevBlockHash
	}
	return currentBlock
}
