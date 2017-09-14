package build

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
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

		// validate and return private key
		keyContents, err := getKey(keyFile)
		if err != nil {
			lumber.Error("hooks:getKey(%s): %s", keyFile, err.Error())
			continue
		}

		// add the keyFile to the list
		keyFiles[file.Name()] = string(keyContents)
	}

	return keyFiles
}

// getKey returns a key's bytes for use in fetching dependencies
func getKey(keyFile string) ([]byte, error) {
	pemBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		// display notice to user failed to read a file
		fmt.Printf("    - Skipping ssh key '%s' (failed to read)\n", keyFile)
		return nil, fmt.Errorf("hooks:ioutil.ReadFile(%s): %s", keyFile, err.Error())
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("No key found")
	}
	buf := block.Bytes

	if encryptedBlock(block) {
		if x509.IsEncryptedPEMBlock(block) {
			// prompt for password to decrypt key
			fmt.Printf("Password protected key found!\nPlease enter the password for '%s'\n", keyFile)
			for attempts := 0; attempts < 3; attempts++ {

				passPhrase, err := display.ReadPassword(keyFile)
				if err != nil {
					return nil, fmt.Errorf("Failed to read password")
				}
				buf, err = x509.DecryptPEMBlock(block, []byte(passPhrase))
				if err != nil {
					if err == x509.IncorrectPasswordError {
						// try again, prompt 3x (mimic `sudo` prompt behavior)
						if attempts == 2 {
							fmt.Println("Too many incorrect password attempts, ignoring key.")
							return nil, err
						}
						fmt.Println("Sorry, try again.")
						continue
					}
					// display notice to user failed to parse a file
					fmt.Printf("    - Skipping ssh key '%s' (failed to decrypt)\n", keyFile)
					return nil, fmt.Errorf("Cannot decode encrypted private keys: %v", err)
				} else {
					break
				}
			}
			block.Headers = map[string]string{}
		} else {
			return nil, fmt.Errorf("Key not encrypted PEM block")
		}
	}

	block.Bytes = buf

	return pem.EncodeToMemory(block), nil
}

// encryptedBlock tells whether a private key is
// encrypted by examining its Proc-Type header
// for a mention of ENCRYPTED
// according to RFC 1421 Section 4.6.1.1.
func encryptedBlock(block *pem.Block) bool {
	return strings.Contains(block.Headers["Proc-Type"], "ENCRYPTED")
}

// if there is just one key specified in the config
// set that as the single key value
func singleKey(keyFile string) map[string]string {
	keyFiles := map[string]string{}

	// validate and return private key
	keyContents, err := getKey(keyFile)
	if err != nil {
		lumber.Error("hooks:getKey(%s): %s", keyFile, err.Error())
		return keyFiles
	}

	keyFiles[filepath.Base(keyFile)] = string(keyContents)
	return keyFiles
}

// isValidKeyFile returns true if a file is a valid key file
func isValidKeyFile(file os.FileInfo) bool {
	return !file.IsDir() &&
		file.Name() != "authorized_keys" &&
		file.Name() != "config" &&
		file.Name() != "known_hosts"
}
