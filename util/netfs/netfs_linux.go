// +build linux

package netfs

import (
	"os/exec"
	"fmt"
	
	"github.com/jcelliott/lumber"
)

// reloadServer reloads the nfs server with the new export configuration
func reloadServer() error {
	// reload nfs server
	//  TODO: provide a clear error message for a direction to fix
	cmd := exec.Command("exportfs", "-ra")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("update: %s", b)
		return fmt.Errorf("update: %s %s", b, err.Error())
	}
	
	return nil
}
