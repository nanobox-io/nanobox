// +build darwin

package netfs

import (
	"os/exec"

	"github.com/jcelliott/lumber"
)

// reloadServer will reload the nfs server with the new export configuration
func reloadServer() error {

	// TODO: make sure nfsd is enabled

	// check the exports to make sure a reload will be successful; TODO: provide a
	// clear message for a direction to fix
	cmd := exec.Command("nfsd", "checkexports")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("output: %s", b)
		return err
	}

	// update exports; TODO: provide a clear error message for a direction to fix
	cmd = exec.Command("nfsd", "update")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("output: %s", b)
		return err
	}

	return nil
}
