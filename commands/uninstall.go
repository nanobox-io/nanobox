// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package commands

//
import (
	"fmt"
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/spf13/cobra"
	"os"
)

//
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls nanobox",
	Long:  ``,

	Run: uninstall,
}

// uninstall
func uninstall(ccmd *cobra.Command, args []string) {

	//
	switch print.Prompt("Are you sure you want to uninstall nanobox (y/N)? ") {

	// don't update by default
	default:
		fmt.Println("Nanobox has not been uninstalled!")
		return

	// if yes continue to update
	case "Yes", "yes", "Y", "y":
		break
	}

	fmt.Println("Uninstalling nanobox... ")

	// do we need to do more here than just this?
	// - shutdown/destroy all vms?
	// - remove virtualbox/vagrant?
	// - probably need to remove nanobox binary

	//
	if err := os.RemoveAll(config.Root); err != nil {
		config.Fatal("[install] os.Remove() failed", err.Error())
	}

	fmt.Println("Nanobox has been successfully uninstalled!")
}
