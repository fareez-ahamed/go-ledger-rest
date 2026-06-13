package main

import (
	"os"

	"github.com/fareez-ahamed/go-ledger-rest/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
