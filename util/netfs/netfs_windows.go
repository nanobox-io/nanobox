// +build windows

package netfs

import (
	"errors"
)

// Add exports a cifs share
func Add(entry string) error {
	return errors.New("Adding an export is not yet supported on windows")
}

// Remove removes a cifs share
func Remove(entry string) error {
	return errors.New("Removing an export is not yet supported on windows")
}

// Exists checks to see if the share already exists
func Exists(entry string) bool {
	return false
}

// Mount mounts a cifs share on a guest machine
func Mount(host_path, mount_path string, context []string) error {
	return errors.New("Mounting a cifs share is not yet supported on windows")
}
