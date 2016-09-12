package fileutil

import (
	"os"
)

// Determines if a file or directory exists
func Exists(path string) bool {

	// stat the file and check for an error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
