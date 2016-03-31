package util

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

	"github.com/nanobox-io/nanobox-boxfile"
)
	
func GlobalDir() string {
	// set Home based off the users homedir (~)
	p, err := homedir.Dir()
	if err != nil {
		// Log.Fatal("[config/config] homedir.Dir() failed", err.Error())
		return ""
	} 
	globalDir := filepath.ToSlash(filepath.Join(p, ".nanobox"))
	os.MkdirAll(globalDir, 0755)
	return globalDir
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

func AppDir() string {
	// if no name is given use localDirName
	app := LocalDirName()

	// read boxfile and look for dev->name
	box := boxfile.NewFromPath(filepath.ToSlash(filepath.Join(GlobalDir(), "boxfile.yml")))
	if box.Valid && (devName := box.Node("dev").StringValue("name"); devName != "") {
		app = devName
	}
	return filepath.ToSlash(filepath.Join(GlobalDir(), "apps", app))
}