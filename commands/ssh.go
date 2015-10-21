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
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
)

//
var sshCmd = &cobra.Command{
	Hidden: true,

	Use:   "ssh",
	Short: "SSH into the nanobox",
	Long:  ``,

	PreRun: boot,
	Run:    ssh,
}

// ssh
func ssh(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("SSHing into nanobox..."))
	if err := Vagrant.SSH(); err != nil {
		Config.Fatal("[commands/ssh] vagrant.SSH() failed - ", err.Error())
	}
}
