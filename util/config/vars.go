package config

import (
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nanobox-io/nanobox-boxfile"
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

	// find out the full path
	cmd := exec.Command("which", programName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// if which doesnt work fall back to just the program name
		return programName
	}

	// trim off any whitespace
	return strings.TrimSpace(string(output))
}
