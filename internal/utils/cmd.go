package utils

import "os"

// StdinHasData checks if standard input has any data to be processed.
func StdinHasData() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}
