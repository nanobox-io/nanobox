//
package server

import (
	"fmt"
	"github.com/nanobox-io/nanobox/config"
	"os"
)

// Suspend
func Suspend() error {

	// if the CLI is running in background mode dont suspend the VM
	if config.VMfile.IsMode("background") {
		fmt.Printf("\n   Note: nanobox is running in background mode. To suspend it run 'nanobox stop'\n\n")
		os.Exit(0)
	}

	//
	res, err := Put("/suspend", nil)
	if err != nil {
		return err
	}

	// only a 2* status code is suspendable
	config.VMfile.SuspendableIs(res.StatusCode/100 == 2)

	// if there are any active consoles dont suspend the VM
	if !config.VMfile.IsSuspendable() {
		fmt.Printf("\n   Note: nanobox has NOT been suspended because there are other active console sessions.\n\n")
		os.Exit(0)
	}

	return nil
}
