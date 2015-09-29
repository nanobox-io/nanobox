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
	Short: "Resumes the suspended nanobox VM",
	Long: `
Description:
  Resumes the halted/suspended nanobox VM by issuing a "vagrant resume"`,

	PreRun: nanoInit,
	Run:    nanoResume,
}

// nanoResume runs 'vagrant resume'
func nanoResume(ccmd *cobra.Command, args []string) {

	// PreRun: nanoInit

	fmt.Printf(stylish.Bullet("Resuming nanobox VM..."))
	if err := util.RunVagrantCommand(exec.Command("vagrant", "resume")); err != nil {
		config.Fatal("[commands/resume] util.RunVagrantCommand() failed", err.Error())
	}
}
