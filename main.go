// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package main

import (
	"fmt"
	"os/exec"

	// api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/commands"
)

// main
func main() {

	var err error

	// ensure vagrant is installed
	err = exec.Command("vagrant").Run()
	if err != nil {
		fmt.Println("Nanobox required 'Vagrant' (https://www.vagrantup.com/) to run. Please download and install it to continue.")
	}

	// ensure virtualbox is installed
	err = exec.Command("vboxmanage").Run()
	if err != nil {
		fmt.Println("Nanobox requires 'Virtualbox' (https://www.virtualbox.org/wiki/Downloads) to run. Please download and install it to continue.")
	}

	if err != nil {
		return
	}

	// do a quick ping to make sure we can communicate properly with the API
	// if err := api.DoRawRequest(nil, "GET", "https://api.nanobox.io/v1/ping", nil, nil); err != nil {
	// 	config.Fatal("[main] The CLI was unable to communicate with the API", err.Error())
	// }

	// check for updates
	// checkUpdate()

	//
	commands.NanoboxCmd.Execute()
}
