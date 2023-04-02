package main

import (
	"flag"
	"fmt"
	blockchain "goChain/blockChain"
	"log"
	"os"
	"runtime"
	"strconv"
)

type Commandline struct {
}

func (cli *Commandline) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getBalance -address ADDRESS - Get Balance For Address")
	fmt.Println(" createBlockChain -address ADDRESS creates a blockChain")
	fmt.Println(" printChain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount - Send Token From To")
}

func (cli *Commandline) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit() // 가바지를 치워버리기 위해서
	}
}

func (cli *Commandline) printAllChain() {
	fmt.Println("Print All Blocks")

	chain := blockchain.ContinueBlockChain("")

	defer chain.Database.Close()
	iter := chain.Iterator()
	block := iter.Next() // 가장 최근 블록을 가져 온다.

	for {
		fmt.Println("\n")
		fmt.Printf("prev Hash  ----> %x \n", block.PrevHash)
		fmt.Printf("Hash  ---->  %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Pow:  ----> %s\n", strconv.FormatBool(pow.Validate()))

		block = iter.GetByPrevHash(block.PrevHash)
		if block == nil {
			return
		}
	}
}

func (cli *Commandline) run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockChainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockChain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockChainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printAllChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

func (cli *Commandline) createBlockChain(address string) {

	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()

	fmt.Println("Chain Is Created")
}

func (cli *Commandline) getBalance(address string) {

	chain := blockchain.ContinueBlockChain(address)

	defer chain.Database.Close()

	balance := 0

	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *Commandline) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)

	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})

	fmt.Println("Success To Send Token")
}

func main() {
	defer os.Exit(0)
	cli := Commandline{}
	cli.run()
}
