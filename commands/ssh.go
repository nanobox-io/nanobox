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
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var sshCmd = &cobra.Command{
	Hidden: true,

	Use:   "ssh",
	Short: "SSH into the nanobox VM",
	Long: `
Description:
  SSHes into the nanobox VM by issuing a "vagrant ssh"`,

	PreRun: bootVM,
	Run:    nanoSSH,
}

// nanoSSH
func nanoSSH(ccmd *cobra.Command, args []string) {

	// PreRun: bootVM

	fmt.Printf(stylish.Bullet("SSHing into nanobox VM..."))

	// NOTE: this command is run manually (vs util.RunVagrantCommand) because the output
	// needs to be hooked up a little different to accomodate Stdin

	// run the command from ~/.nanobox/apps/<this app>
	if err := os.Chdir(config.AppDir); err != nil {
		config.Fatal("[commands/ssh] os.Chdir() failed", err.Error())
	}

	cmd := exec.Command("vagrant", "ssh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//
	fmt.Printf(stylish.Bullet("running '%v'", strings.Trim(fmt.Sprint(cmd.Args), "[]")))

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	if err := cmd.Run(); err != nil {
		config.Fatal("[commands/ssh] cmd.Run() failed", err.Error())
	}
}
