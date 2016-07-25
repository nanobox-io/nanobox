// +build windows

package netfs

import (
	"errors"
)


// https://gist.github.com/notxarb/ebe664d693710ddb2b110d1251800e67

// maybe a way to mount without the password:
// 	https://ubuntuforums.org/showthread.php?t=1589669

// Exists checks to see if the share already exists
func Exists(entry string) bool {
	return false
}

// Add exports a cifs share
func Add(entry string) error {
	return errors.New("Adding an export is not yet supported on windows")
}

// Remove removes a cifs share
func Remove(entry string) error {
	return errors.New("Removing an export is not yet supported on windows")
}

// Mount mounts a cifs share on a guest machine
func Mount(host_path, mount_path string) error {
	return errors.New("Mounting a cifs share is not yet supported on windows")
}
