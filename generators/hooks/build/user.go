package build

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jcelliott/lumber"
	"golang.org/x/crypto/ssh"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
)

// UserPayload returns a string for the user hook payload
func UserPayload() string {

	configModel, _ := models.LoadConfig()
	payload := map[string]interface{}{
		"provider": configModel.Provider,
		"os":       runtime.GOOS,
	}

	payload["ssh_files"] = sshKeys()

	// marshal the payload into json
	b, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}

	return string(b)
}

// collect the ssh keys from the specified location and return them as a map
func sshKeys() map[string]string {
	keyFiles := map[string]string{}
	sshFolder := config.SSHDir()
	configModel, _ := models.LoadConfig()
	if configModel.SshKey != "default" {
		fi, err := os.Stat(configModel.SshKey)
		// if i were able to read the file then we will be altering the ssh key collection
		if err == nil {
			// if it isnt a directory the keys will be a single key pair
			if !fi.IsDir() {
				return singleKey(configModel.SshKey)
			}
			// if it is a directory continue as normal but using this directory instead
			sshFolder = configModel.SshKey
		}
	}

	// read all of the ssh files on the system
	sshFiles, err := ioutil.ReadDir(sshFolder)
	if err != nil {
		return keyFiles
	}

	for _, file := range sshFiles {

		if !isValidKeyFile(file) {
			continue
		}

		// read the contents of the keyfile
		keyFile := filepath.Join(sshFolder, file.Name())
		content, err := ioutil.ReadFile(keyFile)
		// todo: display notice to user failed to read a file
		if err != nil {
			lumber.Error("hooks:ioutil.ReadFile(%s): %s", keyFile, err.Error())
			// if this file cant be read continue on and give as many
			// of the ssh keys as we can
			continue
		}

		// ensure key is not password protected and is a valid private key file
		_, err = ssh.ParsePrivateKey(content)
		if err != nil {
			lumber.Error("hooks:ssh.ParsePrivateKey(%s): %s", keyFile, err.Error())
			continue
		}

		// add the keyFile to the list
		keyFiles[file.Name()] = string(content)
	}

	return keyFiles
}

// if there is just one key specified in the config
// set that as the single key value
func singleKey(key string) map[string]string {
	keyFiles := map[string]string{}
	content, err := ioutil.ReadFile(key)
	if err != nil {
		return keyFiles
	}
	keyFiles[filepath.Base(key)] = string(content)
	return keyFiles
}

// isValidKeyFile returns true if a file is a valid key file
func isValidKeyFile(file os.FileInfo) bool {
	return !file.IsDir() &&
		file.Name() != "authorized_keys" &&
		file.Name() != "config" &&
		file.Name() != "known_hosts"
}
