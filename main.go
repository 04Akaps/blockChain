package main

import (
	"github.com/jjimgo/blockChain.git/wallet"
	"os"
)

func main() {
	defer os.Exit(0)
	//cli := cli.Commandline{}
	//cli.Run()

	w := wallet.MakeWallet()
	w.Address()
}
