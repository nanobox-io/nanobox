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
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var resumeCmd = &cobra.Command{
	Hidden: true,

	Use:   "resume",
	Short: "Resumes the nanobox",
	Long:  ``,

	PreRun: nanoInit,
	Run:    nanoResume,
}

// nanoResume runs 'vagrant resume'
func nanoResume(ccmd *cobra.Command, args []string) {

	// PreRun: nanoInit

	fmt.Printf(stylish.Bullet("Resuming nanobox..."))
	if err := util.VagrantRun(exec.Command("vagrant", "resume")); err != nil {
		config.Fatal("[commands/resume] util.VagrantRun() failed", err.Error())
	}
}
