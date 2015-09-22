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
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var createCmd = &cobra.Command{
	Hidden: true,

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

	// if the command is being run with the "add" flag, it means an entry needs to
	// be added to the hosts file and execution yielded back to the parent
	if len(args) > 0 && args[0] == "add" {
		util.AddDevDomain()
		os.Exit(0)
	}

	// run an init to ensure there is a Vagrantfile
	nanoInit(nil, args)

	//
	// boot the vm
	fmt.Printf(stylish.ProcessStart("creating nanobox vm"))
	if err := util.RunVagrantCommand(exec.Command("vagrant", "up")); err != nil {
		util.Fatal("[commands/create] util.RunVagrantCommand() failed", err)
	}

	//
	// open the /etc/hosts file for scanning...
	f, err := os.Open("/etc/hosts")
	if err != nil {
		util.Fatal("[commands/create] os.Open() failed", err)
	}
	defer f.Close()

	// determines whether or not an entry needs to be added to the /etc/hosts file
	// (an entry will be added unless it's confirmed that it's not needed)
	addEntry := true

	// scan hosts file looking for an entry corresponding to this app...
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// if an entry with the IP is detected, flag the entry as not needed
		if strings.HasPrefix(scanner.Text(), config.Nanofile.IP) {
			addEntry = false
		}
	}

	// add the entry if needed
	// if addEntry && util.AccessDenied() {
	if addEntry {
		util.SudoExec("create add", "Attempting to add nano.dev domain to hosts file")
	}

	// if devmode is detected, the machine needs to be rebooted for devmode to take
	// effect
	if fDevmode {
		fmt.Printf(stylish.Bullet("Rebooting machine to finalize 'devmode' configuration..."))
		nanoReload(nil, args)
	}
}
