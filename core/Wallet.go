package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
)

const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey `json:"privateKey"`
	PublicKey  []byte           `json:"publicKey"`
}

func NewWallet() *Wallet {
	curve := elliptic.P256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	w := &Wallet{*private, pubKey}
	return w
}

func (w *Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	checksum := checksum(pubKeyHash)
	fullPayload := append(pubKeyHash, checksum...)
	address := Base58Encode(fullPayload)
	return address
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(publicSHA256[:])
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}

func (w *Wallet) Serialize() []byte {
	str, _ := json.Marshal(w)
	return str
}

func DeserializeWallet(js []byte) *Wallet {
	var w Wallet
	json.Unmarshal(js, &w)
	return &w
}
