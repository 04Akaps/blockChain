package main

import (
	"github.com/jjimgo/blockChain.git/cli"
	"os"
)

func main() {
	defer os.Exit(0)
	cli := cli.Commandline{}
	cli.Run()

}
