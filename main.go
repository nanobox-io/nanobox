// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package main

import (
	"fmt"
	"os/exec"

	"github.com/nanobox-io/nanobox-cli/commands"
)

// main
func main() {

	pass := true

	// ensure vagrant is installed
	if err := exec.Command("vagrant", "-v").Run(); err != nil {
		fmt.Println("Missing dependency 'Vagrant'. Please download and install it to continue (https://www.vagrantup.com/).")
		pass = false
	}

	// ensure virtualbox is installed
	if err := exec.Command("vboxmanage", "-v").Run(); err != nil {
		fmt.Println("Missing dependency 'Virtualbox'. Please download and install it to continue (https://www.virtualbox.org/wiki/Downloads).")
		pass = false
	}

	// if a dependency check fails, exit
	if !pass {
		return
	}

	// check for updates
	// checkUpdate()

	//
	commands.NanoboxCmd.Execute()
}
