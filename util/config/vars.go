package config

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

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

// AppID is the id of the app we will use based
// on the name as well as the folder
// this should help us keep a unique name for apps that 
// happen to have the same folder base name
func AppID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(LocalDir())))
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
			if content, err := ioutil.ReadFile(filepath.Join(SSHDir(), file.Name())); err == nil {
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
