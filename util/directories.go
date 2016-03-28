package util

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func GlobalDir() string {
	// set Home based off the users homedir (~)
	p, err := homedir.Dir()
	if err != nil {
		// Log.Fatal("[config/config] homedir.Dir() failed", err.Error())
		return ""
	} 
	return filepath.ToSlash(filepath.Join(p, ".nanobox"))
}

func LocalDir() string {
	p, err := os.Getwd()
	if err != nil {
		// Log.Fatal("[config/config] os.Getwd() failed", err.Error())
		return ""
	}
	return filepath.ToSlash(p)
}

func LocalDirName() string {
	return filepath.Base(LocalDir())
}