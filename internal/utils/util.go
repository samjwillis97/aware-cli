package utils

import (
	"fmt"
	"os"
)

func ExitIfError(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(1)
}
