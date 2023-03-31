package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	block "goChain/blockChain"
	blockchain "goChain/blockChain"
)

type Commandline struct {
	blockchain *blockchain.BlockChain
}

func (cli *Commandline) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block Block_data - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *Commandline) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit() // 가바지를 치워버리기 위해서
	}
}

func (cli *Commandline) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *Commandline) printChain() {
	iter := cli.blockchain.Iterator()

	block := iter.Next()

	fmt.Println("\n")
	fmt.Printf("prev Hash %x\n", block.PrevHash)
	fmt.Printf("Data In Block %s\n", block.Data)
	fmt.Printf("Hash  %x\n", block.Hash)

	pow := blockchain.NewProof(block)
	fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
	fmt.Println("\n")

	//if len(block.PrevHash) == 0 {
	//	break
	//}

}

func (cli *Commandline) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.ErrorHandle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.ErrorHandle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}

		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		// print를 입력했을 떄
		cli.printChain()
	}
}

func main() {
	chain := block.InitBlockChain()
	defer chain.Database.Close()

	cli := Commandline{chain}
	cli.run()
}
