package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jcelliott/lumber"
	"github.com/mitchellh/go-homedir"

	"github.com/nanobox-io/nanobox-boxfile"
)

// TODO: win: make sure the folders are in the correct place
// the go-homedir is supposed to work in both packages
// and we used filpath.ToSlash which puts the correct slashes in.

// GlobalDir ...
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

// SSHDir ...
func SSHDir() string {
	p, err := homedir.Dir()
	if err != nil {
		// Log.Fatal("[config/config] homedir.Dir() failed", err.Error())
		return ""
	}
	return filepath.ToSlash(filepath.Join(p, ".ssh"))
}

// LocalDir ...
func LocalDir() string {
	p, err := os.Getwd()
	if err != nil {
		// Log.Fatal("[config/config] os.Getwd() failed", err.Error())
		return ""
	}
	return filepath.ToSlash(p)
}

// LocalDirName ...
func LocalDirName() string {
	return filepath.Base(LocalDir())
}

// BoxfileLocation ...
func BoxfileLocation() string {
	return filepath.ToSlash(filepath.Join(LocalDir(), "boxfile.yml"))
}

// AppName ...
func AppName() string {
	// if no name is given use localDirName
	app := LocalDirName()

	// read boxfile and look for dev->name
	box := boxfile.NewFromPath(BoxfileLocation())
	devName := box.Node("dev").StringValue("name")
	if devName != "" {
		app = devName
	}
	return app
}

// EngineDir gets the directory of the engine if it is a directory and on the
// local file system
func EngineDir() string {
	box := boxfile.NewFromPath(BoxfileLocation())
	engineName := box.Node("env").StringValue("engine")
	if engineName != "" {
		fi, err := os.Stat(engineName)
		if err == nil && fi.IsDir() {
			return engineName
		}
	}
	return ""
}

// UserPayload ...
func UserPayload() string {
	sshFiles, err := ioutil.ReadDir(SSHDir())
	if err != nil {
		fmt.Println("upserpayload", err)
		return `{"ssh_files":{}}`
	}
	files := map[string]string{}
	for _, file := range sshFiles {
		if !file.IsDir() && file.Name() != "authorized_keys" && file.Name() != "config" && file.Name() != "known_hosts" {
			content, err := ioutil.ReadFile(filepath.Join(SSHDir(), file.Name()))
			if err == nil {
				files[file.Name()] = string(content)
			}
		}
	}
	b, err := json.Marshal(map[string]interface{}{"ssh_files": files})
	if err != nil {
		fmt.Println("upserpayload", err)
		return `{"ssh_files":{}}`
	}
	return string(b)
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

	//
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

	//
	return
}
