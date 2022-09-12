package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

func Success(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, fmt.Sprintf("\n\u001B[0;32m✓\u001B[0m %s\n", msg), args...)
}

func ShowLoading(msg string) *spinner.Spinner {
	const refreshRate = 100 * time.Millisecond

	s := spinner.New(
		spinner.CharSets[14],
		refreshRate,
		spinner.WithSuffix(" "+msg),
		spinner.WithHiddenCursor(true),
		spinner.WithWriter(color.Error),
	)
	s.Start()

	return s
}

func Fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", msg), args...)
}

func Failed(msg string, args ...interface{}) {
	Fail(msg, args...)
	os.Exit(1)
}
