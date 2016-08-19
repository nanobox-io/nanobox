package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcelliott/lumber"
	"github.com/mitchellh/go-homedir"

	"github.com/nanobox-io/nanobox-boxfile"
)

// GlobalDir ...
func GlobalDir() string {

	// set Home based off the users homedir (~)
	p, err := homedir.Dir()
	if err != nil {
		return ""
	}

	globalDir := filepath.ToSlash(filepath.Join(p, ".nanobox"))
	os.MkdirAll(globalDir, 0755)

	return globalDir
}

// ImplodeGlobalDir will remove the global dir and everything inside
func ImplodeGlobalDir() error {
	globalDir := GlobalDir()
	
	// blast the global directory
	if err := os.RemoveAll(globalDir); err != nil {
		return fmt.Errorf("failed to remove %s: %s", globalDir, err.Error())
	}
	
	return nil
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

	box := boxfile.NewFromPath(Boxfile())
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
