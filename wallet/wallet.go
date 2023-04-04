package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	blockchain "github.com/jjimgo/blockChain.git/blockChain"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	blockchain.ErrorHandle(err)

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()

	return &Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}
}
