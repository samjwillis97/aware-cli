package utils

import (
	"fmt"
	"os"
	"ampaware.com/cli/pkg/aware"
)

func ExitIfError(err error) {
    if err == nil {
        return
    }

    var msg string

    switch err {
    case aware.ErrEmptyResult:
        msg = "aware: Received empty response.\n Please try again."
    default:
        msg = fmt.Sprintf("Error: %s", err.Error())
    }

    fmt.Fprintf(os.Stderr, "%s\n", msg)
    os.Exit(1)
}
