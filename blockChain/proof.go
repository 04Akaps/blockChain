package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 12

// 채굴 어려움을 의미
// 원래는 알고리즘에 의해서 달라지는 값이어야 함

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	target.Lsh(target, uint(256-Difficulty))
	// Lsh는 target를 target=x<=n으로 만들고 target을 반환

	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer) // 단순히 메모리를 할당하기 위해서
	// make와 &bytes.Buffer 이런식으로도 사용이 가능하다.
	// make는 값을 초기화 까지 할 떄 사용하고
	// pointer로 직접 선언하는 것은 Effective Go에 따르면 차이는 없다고 한다.
	err := binary.Write(buff, binary.BigEndian, num)
	// buff에 값을 쓰기 위한 행위
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
