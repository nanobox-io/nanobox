package config

import (
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"

	"github.com/nanobox-io/nanobox-boxfile"
	
	"github.com/nanobox-io/nanobox/util/fileutil"
)

// AppName ...
func AppName() string {

	// if no name is given use localDirName
	app := LocalDirName()

	// read boxfile and look for dev:name
	box := boxfile.NewFromPath(Boxfile())
	devName := box.Node("dev").StringValue("name")

	// set the app name
	if devName != "" {
		app = devName
	}

	return app
}

// EnvID ...
func EnvID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(LocalDir())))
}

// NanoboxPath ...
func NanoboxPath() string {
	
	programName := os.Args[0]
	
	// if args[0] was a path to nanobox already
	if fileutil.Exists(programName) {
		return programName
	}
	
	// lookup the full path to nanobox
	path, err := exec.LookPath(programName)
	if err == nil {
		return path
	}

	// unable to find the full path, just return what was called
	return programName
}
