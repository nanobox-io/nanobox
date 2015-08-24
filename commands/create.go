// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Runs 'nanobox init', then boots the nanobox VM",
	Long: `
Description:
  Runs 'nanobox init', then boots the nanobox VM by issuing a "vagrant up"`,

	Run: nanoCreate,
}

//
// nanoCreate
func nanoCreate(ccmd *cobra.Command, args []string) {

	//
	// open the /etc/hosts file for scanning...
	f, err := os.Open("/etc/hosts")
	if err != nil {
		ui.LogFatal("[commands/create] os.Open() failed", err)
	}
	defer f.Close()

	// a new scanner for scanning the /etc/hosts file
	scanner := bufio.NewScanner(f)

	// determines whether or not an entry needs to be added to the /etc/hosts file
	// (an entry will be added unless it's confirmed that it's not needed)
	addEntry := true

	// scan hosts file looking for an entry corresponding to this app...
	for scanner.Scan() {

		// if an entry with the IP is detected, flag the entry as not needed
		if strings.HasPrefix(scanner.Text(), config.Nanofile.IP) {
			addEntry = false
		}
	}

	// add the entry if needed
	if addEntry {
		if utils.AccessDenied() {
			utils.SudoExec("create", "Attempting to add nano.dev domain to hosts file")
			os.Exit(0)
		}

		utils.AddDevDomain()
	}

	// run an init to ensure there is a Vagrantfile
	nanoInit(nil, args)

	//
	// boot the vm
	fmt.Printf(stylish.ProcessStart("starting nanobox vm"))
	if err := runVagrantCommand(exec.Command("vagrant", "up")); err != nil {
		ui.LogFatal("[commands/create] runVagrantCommand() failed", err)
	}
	fmt.Printf(stylish.ProcessEnd())

	// upgrade all nanobox docker images
	nanoUpgrade(nil, args)
}

//
func create() {
	utils.AddDevDomain()
}
