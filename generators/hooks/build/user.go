package build

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/config"
)

// UserPayload returns a string for the user hook payload
func UserPayload() string {

	// read all of the ssh files on the system
	sshFiles, err := ioutil.ReadDir(config.SSHDir())
	if err != nil {
		return "{}"
	}

	// create a list of files that will be installed on the system
	keyFiles := map[string]string{}
	for _, file := range sshFiles {

		if !isValidKeyFile(file) {
			continue
		}

		// read the contents of the keyfile
		keyFile := filepath.Join(config.SSHDir(), file.Name())
		content, err := ioutil.ReadFile(keyFile)
		if err != nil {
			lumber.Error("hooks:ioutil.ReadFile(%s): %s", keyFile, err.Error())
			// if this file cant be read continue on and give as many
			// of the ssh keys as we can
			continue
		}

		// add the keyFile to the list
		keyFiles[file.Name()] = string(content)
	}

	payload := map[string]interface{}{
		"ssh_files": keyFiles,
	}

	// marshal the payload into json
	b, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}

	return string(b)
}

// isValidKeyFile returns true if a file is a valid key file
func isValidKeyFile(file os.FileInfo) bool {

	return !file.IsDir() &&
		file.Name() != "authorized_keys" &&
		file.Name() != "config" &&
		file.Name() != "known_hosts"
}
