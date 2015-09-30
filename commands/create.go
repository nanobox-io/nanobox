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

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var createCmd = &cobra.Command{
	Hidden: true,

	Use:   "create",
	Short: "Creates a new nanobox",
	Long:  ``,

	PreRun: nanoInit,
	Run:    nanoCreate,
}

//
func init() {
	createCmd.Flags().BoolVarP(&fAddEntry, "add-entry", "", false, "")
	createCmd.Flags().MarkHidden("add-entry")
}

//
// nanoCreate
func nanoCreate(ccmd *cobra.Command, args []string) {

	// PreRun: nanoInit

	// if the command is being run with the "add" flag, it means an entry needs to
	// be added to the hosts file and execution yielded back to the parent
	if fAddEntry {
		util.HostfileAddDomain()
		os.Exit(0)
	}

	//
	// boot the vm
	fmt.Printf(stylish.Bullet("Creating a nanobox"))
	if err := util.VagrantRun(exec.Command("vagrant", "up")); err != nil {
		config.Fatal("[commands/create] util.VagrantRun() failed", err.Error())
	}

	// after the machine boots, update the docker images
	imagesUpdate(nil, args)

	//
	// open the /etc/hosts file for scanning...
	f, err := os.Open("/etc/hosts")
	if err != nil {
		config.Fatal("[commands/create] os.Open() failed", err.Error())
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
		util.SudoExec("create --add-entry", fmt.Sprintf("Adding %v domain to hosts file", config.Nanofile.Domain))
	}

	// if devmode is detected, the machine needs to be rebooted for devmode to take
	// effect
	if fDevmode {
		fmt.Printf(stylish.Bullet("Rebooting machine to finalize 'devmode' configuration..."))
		nanoReload(nil, args)
	}
}
