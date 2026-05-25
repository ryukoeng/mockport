package main

import (
	"os"

	"github.com/albert-einshutoin/mockport/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
