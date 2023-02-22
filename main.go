package main

import (
	"fmt"
	"strconv"

	block "goChain/blockChain"
	blockchain "goChain/blockChain"
)

func main() {
	chain := block.InitBlockChain()

	chain.AddBlock("first Block")
	chain.AddBlock("second Block")
	chain.AddBlock("third Block")

	for _, block := range chain.Blocks {
		fmt.Println("\n")
		fmt.Printf("prev Hash %x\n", block.PrevHash)
		fmt.Printf("Data In Block %s\n", block.Data)
		fmt.Printf("Hash  %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))

	}
}
