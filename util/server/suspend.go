// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package server

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/config"
)

// Suspend
func Suspend() error {

	// if the CLI is running in background mode dont suspend the VM
	if config.VMfile.IsMode("background") {
		fmt.Printf("\n   Note: nanobox is running in background mode. To suspend it run 'nanobox down'\n\n")
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
