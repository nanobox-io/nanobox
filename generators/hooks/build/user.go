package build

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  
  "github.com/nanobox-io/nanobox/util/config"
)

// UserPayload returns a string for the user hook payload
func UserPayload() (string, error) {
  
  // read all of the ssh files on the system
  sshFiles, err := ioutil.ReadDir(config.SSHDir())
  if err != nil {
    return "", fmt.Errorf("failed to read ssh directory: %s", err.Error())
  }
  
  // create a list of files that will be installed on the system
  keyFiles := map[string]string{}
  for _, file := range sshFiles {
    
    if !isValidKeyFile(file) {
      continue
    }
    
    // read the contents of the keyfile
    keyFile := filepath.Join(config.SSHDir(), file.Name())
    content, err := ioutil.ReadFile(keyFile); 
    if err != nil {
      return "", fmt.Errorf("failed to read ssh key file (%s): %s", keyFile, err.Error())
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
    return "", fmt.Errorf("failed to encode hook payload into json: %s", err.Error())
  }
  
  return string(b), nil
}

// isValidKeyFile returns true if a file is a valid key file
func isValidKeyFile(file os.FileInfo) bool {
  // ignore directories
  if file.IsDir() {
    return false
  }
  
  // ignore authorized_keys file
  if file.Name() == "authorized_keys" {
    return false
  }
  
  // ignore ssh config
  if file.Name() == "config" {
    return false
  }
  
  // ignore known_hosts
  if file.Name() == "known_hosts" {
    return false
  }
  
  return true
}
