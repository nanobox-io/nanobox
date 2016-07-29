// +build darwin

package netfs

import (
	"fmt"
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
		lumber.Debug("checkexports: %s", b)
		return fmt.Errorf("checkexports: %s %s", b, err.Error())
	}

	// update exports; TODO: provide a clear error message for a direction to fix
	cmd = exec.Command("nfsd", "update")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("update: %s", b)
		return fmt.Errorf("update: %s %s", b, err.Error())
	}

	return nil
}
