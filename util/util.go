// Package util ...
package util

import (
	"os"
	"math/rand"
	"time"
)

const (

	// VERSION is the global version for nanobox; mainly used in the update process
	// but placed here to allow access when needed (commands, processor, etc.)
	VERSION     = "1.0.0"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// RandomString ...
func RandomString(size int) string {

	// create a new randomizer with a unique seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	//
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}

	return string(b)
}

func FolderExists(folderName string) bool {
	dir, err := os.Stat(folderName)
	if err != nil {
		return false
	}
	return dir.IsDir()
}
