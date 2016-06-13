package util

import (
	"os"
	"path/filepath"

	"github.com/jcelliott/lumber"
	"github.com/mitchellh/go-homedir"

	"github.com/nanobox-io/nanobox-boxfile"
)

// TODO: win: make sure the folders are in the correct place the go-homedir is
// supposed to work in both packages and we used filpath.ToSlash which puts the
// correct slashes in.

// GlobalDir ...
func GlobalDir() string {

	// set Home based off the users homedir (~)
	p, err := homedir.Dir()
	if err != nil {
		return ""
	}

	//
	globalDir := filepath.ToSlash(filepath.Join(p, ".nanobox"))
	os.MkdirAll(globalDir, 0755)

	return globalDir
}

// LocalDir ...
func LocalDir() string {

	//
	p, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.ToSlash(p)
}

// LocalDirName ...
func LocalDirName() string {
	return filepath.Base(LocalDir())
}

// SSHDir ...
func SSHDir() string {

	//
	p, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return filepath.ToSlash(filepath.Join(p, ".ssh"))
}

// EngineDir gets the directory of the engine if it is a directory and on the
// local file system
func EngineDir() string {

	box := boxfile.NewFromPath(BoxfileLocation())
	engineName := box.Node("env").StringValue("engine")

	//
	if engineName != "" {
		fi, err := os.Stat(engineName)
		if err == nil && fi.IsDir() {
			return engineName
		}
	}

	return ""
}

// TmpDir creates a temporary directory where nanobox specific files can be
// downloaded (updated versions of nanobox, md5's, etc.)
func TmpDir() (tmpDir string) {

	//
	tmpDir = filepath.ToSlash(filepath.Join(GlobalDir(), "tmp"))

	//
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		lumber.Fatal("[config/config] os.Mkdir() failed", err.Error())
	}

	return
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
