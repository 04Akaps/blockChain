package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct { // Block을 구성할 struct
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
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
