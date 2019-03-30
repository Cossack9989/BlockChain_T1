package core

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

const targetBits = 24
const maxnonce = math.MaxInt64

type Block struct {
	Timestamp     int64  `json:"timestamp"`
	Data          []byte `json:"data"`
	PrevBlockHash []byte `json:"prevBlockHash"`
	Proof         int64  `json:"proof"`
}

func NewGenesisBlock() *Block {
	b := &Block{time.Now().Unix(), []byte("Hello World!"), []byte{}, int64(0)}
	b.ProofOfWork()
	return b
}

func (b *Block) GetBlockHash() []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(strconv.FormatInt(b.Timestamp, 16)),
			b.Data,
			b.PrevBlockHash,
			[]byte(strconv.FormatInt(b.Proof, 16)),
		},
		[]byte{},
	)
	hash := sha256.Sum256(data)
	return hash[:]
}

func (b *Block) ProofOfWork() {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)
	var hashInt big.Int
	var hash [32]byte
	nonce, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	fmt.Printf("Mining the block containing \"%s\"\n", b.Data)
	for nonce.Int64() < maxnonce {
		data := bytes.Join(
			[][]byte{
				[]byte(strconv.FormatInt(b.Timestamp, 16)),
				b.Data,
				b.PrevBlockHash,
				[]byte(strconv.FormatInt(nonce.Int64(), 16)),
			},
			[]byte{},
		)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce, _ = rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		}
	}
	b.Proof = nonce.Int64()
}

func (b *Block) Serialize() []byte {
	result, _ := json.Marshal(b)
	return []byte(result)
}

func DeserializeBlock(js []byte) *Block {
	var b Block
	json.Unmarshal(js, &b)
	return &b
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	b := &Block{time.Now().Unix(), []byte(data), prevBlockHash, int64(0)}
	b.ProofOfWork()
	return b
}
