// +build windows

package netfs

import (
	"errors"
	
	// "github.com/nanobox-io/nanobox/util/config"
)


// https://gist.github.com/notxarb/ebe664d693710ddb2b110d1251800e67

// Exists checks to see if the share already exists
func Exists(path string) bool {

	return false
}

// Add exports a cifs share
func Add(path string) error {
	return errors.New("Adding an export is not yet supported on windows")
}

// Remove removes a cifs share
func Remove(path string) error {
	return errors.New("Removing an export is not yet supported on windows")
}

// Mount mounts a cifs share on a guest machine
func Mount(hostPath, mountPath string) error {
	return errors.New("Mounting a cifs share is not yet supported on windows")
}
