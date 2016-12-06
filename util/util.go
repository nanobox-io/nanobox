// Package util ...
package util

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	"path/filepath"
	
	"github.com/nanobox-io/nanobox/util/config"
)

const (

	// VERSION is the global version for nanobox; mainly used in the update process
	// but placed here to allow access when needed (commands, processor, etc.)
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

func FileMD5(name string) string {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		// give the relative path a chance
		// but if it doesnt attach the filename given to the absolute path
		data, err = ioutil.ReadFile(filepath.ToSlash(filepath.Join(config.LocalDir(), name)))
		if err != nil {
			return ""
		}
	}
	return fmt.Sprintf("%x", md5.Sum(data))
}
