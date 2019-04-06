package core

import (
	"github.com/boltdb/bolt"
)

const walletdb = "Wallets.db"
const walletbucket = "Wallets"

type Wallets struct {
	db *bolt.DB
}

func NewWallets() *Wallets {
	db, _ := bolt.Open(walletdb, 0600, nil)
	ws := &Wallets{db}
	return ws
}

func (ws *Wallets) GetWallet(address string) *Wallet {
	var w *Wallet
	ws.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(walletbucket))
		js := b.Get([]byte(address))
		w = DeserializeWallet(js)
		return nil
	})
	return w
}

func (ws *Wallets) CreateWallet() string {
	w := NewWallet()
	address := w.GetAddress()
	ws.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(walletbucket))
		if b == nil {
			b, _ := tx.CreateBucket([]byte(walletbucket))
			b.Put(address, w.Serialize())
		} else {
			b.Put(address, w.Serialize())
		}
		return nil
	})
	return string(address)
}
