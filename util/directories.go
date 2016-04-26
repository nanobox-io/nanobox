package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func SshDir() string {
	p, err := homedir.Dir()
	if err != nil {
		// Log.Fatal("[config/config] homedir.Dir() failed", err.Error())
		return ""
	}
	return filepath.ToSlash(filepath.Join(p, ".ssh"))
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

func BoxfileLocation() string {
	return filepath.ToSlash(filepath.Join(GlobalDir(), "boxfile.yml"))
}

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

// get the director of the engine if it is a directory
// and on my local file system
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

func UserPayload() string {
	sshFiles, err := ioutil.ReadDir(SshDir())
	if err != nil {
		fmt.Println("upserpayload", err)
		return `{"ssh_files":{}}`
	}
	files := map[string]string{}
	for _, file := range sshFiles {
		if !file.IsDir() && file.Name() != "authorized_keys" && file.Name() != "config" && file.Name() != "known_hosts" {
			content, err := ioutil.ReadFile(filepath.Join(SshDir(), file.Name()))
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
