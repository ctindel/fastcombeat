package main

import (
	"os"

	"github.com/ctindel/fastcombeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
