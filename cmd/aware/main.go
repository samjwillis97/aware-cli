// Package main is the main command entry-point to aware CLI.
package main

import (
	"fmt"
	"os"

	"ampaware.com/cli/internal/cmd/root"
)

func main() {
	rootCmd := root.NewCmdRoot()
	if _, err := rootCmd.ExecuteC(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
