// Package util ...
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"

	"github.com/nanobox-io/nanobox-boxfile"
)

const (

	// VERSION is the global version for nanobox; mainly used in the update process
	// but placed here to allow access when needed (commands, processor, etc.)
	VERSION     = "1.0.0"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// BoxfileLocation ...
func BoxfileLocation() string {
	return filepath.ToSlash(filepath.Join(LocalDir(), "boxfile.yml"))
}

// AppName ...
func AppName() string {

	// if no name is given use localDirName
	app := LocalDirName()

	// read boxfile and look for dev:name
	box := boxfile.NewFromPath(BoxfileLocation())
	devName := box.Node("dev").StringValue("name")

	// set the app name
	if devName != "" {
		app = devName
	}

	return app
}

// UserPayload ...
func UserPayload() string {

	//
	sshFiles, err := ioutil.ReadDir(SSHDir())
	if err != nil {
		fmt.Println("upserpayload", err)
		return `{"ssh_files":{}}`
	}

	//
	files := map[string]string{}
	for _, file := range sshFiles {
		if !file.IsDir() && file.Name() != "authorized_keys" && file.Name() != "config" && file.Name() != "known_hosts" {
			if content, err := ioutil.ReadFile(filepath.Join(SSHDir(), file.Name())); err != nil {
				files[file.Name()] = string(content)
			}
		}
	}

	//
	b, err := json.Marshal(map[string]interface{}{"ssh_files": files})
	if err != nil {
		fmt.Println("upserpayload", err)
		return `{"ssh_files":{}}`
	}

	return string(b)
}

// RandomString ...
func RandomString(size int) string {

	//
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
