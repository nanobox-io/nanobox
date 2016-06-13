// +build linux

package netfs

import (
	"errors"
)

// reloadServer reloads the nfs server with the new export configuration
func reloadServer() error {

	// TODO: figure out how to do this :/
	return errors.New("Reloading an NFS server is not yet implemented on linux")
}
