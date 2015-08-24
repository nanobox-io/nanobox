// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var haltCmd = &cobra.Command{
	Use:   "halt",
	Short: "Halts the running nanobox VM",
	Long: `
Description:
  Halts the running nanobox VM by issuing a "vagrant halt"`,

	Run: nanoHalt,
}

//
func init() {
	haltCmd.Flags().BoolVarP(&fForce, "force", "f", false, "Skips confirmation and forces the nanobox VM to halt")
}

// nanoHalt
func nanoHalt(ccmd *cobra.Command, args []string) {

	if !fForce {

		// prompt for confirmation...
		switch ui.Prompt("Are you sure you want to halt this VM (y/N)? ") {

		// if positive confirmation, proceed and halt
		case "Y", "y":
			fmt.Printf(stylish.Bullet("Halt confirmed, continuing..."))

		// if negative confirmation, exit w/o halting
		default:
			os.Exit(0)
		}
	}

	// halt the vm...
	fmt.Printf(stylish.ProcessStart("halting nanobox vm"))
	if err := runVagrantCommand(exec.Command("vagrant", "halt", "--force")); err != nil {
		ui.LogFatal("[commands/halt] runVagrantCommand() failed", err)
	}
	fmt.Printf(stylish.ProcessEnd())
}
