package wallet

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	blockchain "github.com/jjimgo/blockChain.git/blockChain"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256() // 타원곡선을 통한 256바이트의 값을 출력

	private, err := ecdsa.GenerateKey(curve, rand.Reader) // 해당 값을 통해서 privateKey를 생성

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

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := crypto.RIPEMD160.New()

	_, err := hasher.Write(pubHash[:])

	blockchain.ErrorHandle(err)

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	
	return secondHash[:checksumLength]
}
