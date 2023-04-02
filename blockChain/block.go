package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct { // Block을 구성할 struct
	Hash        []byte
	Transaction []*Transaction
	PrevHash    []byte
	Nonce       int
}

func (b *Block) HashTransactions() []byte {
	// 블록에 있는 데이터를 해시화
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transaction {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func CreateBlock(transactions []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, transactions, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis(genesisTx *Transaction) *Block {
	return CreateBlock([]*Transaction{genesisTx}, []byte{})
}

func (b *Block) Serialize() []byte {
	// block을 인코딩
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	ErrorHandle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	// 특정 데이터를 block에 디코딩
	var block Block

	decoder := gob.NewDecoder(bytes.NewBuffer(data))

	err := decoder.Decode(&block)

	ErrorHandle(err)

	return &block
}

func ErrorHandle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
