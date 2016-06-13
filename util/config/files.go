package config

import (
	"os"
	"path/filepath"

	"github.com/jcelliott/lumber"
)

// Boxfile ...
func Boxfile() string {
	return filepath.ToSlash(filepath.Join(LocalDir(), "boxfile.yml"))
}

// UpdateFile creates an update file thats used in the update process to determine
// when the last time nanobox was updated
func UpdateFile() (updateFile string) {

	//
	updateFile = filepath.ToSlash(filepath.Join(GlobalDir(), ".update"))

	// return the filepath if it's already created...
	if _, err := os.Stat(updateFile); err == nil {
		return
	}

	// ...otherwise create the file
	f, err := os.Create(updateFile)
	if err != nil {
		lumber.Fatal("[config/config] os.Create() failed", err.Error())
	}
	defer f.Close()

	return
}
