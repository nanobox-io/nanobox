package fileutil

import (
	"os"
)

// Determines if a file or directory exists
func Exists(path string) bool {
	// stat the file and check for an error
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}
