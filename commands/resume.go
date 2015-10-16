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

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/vagrant"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var resumeCmd = &cobra.Command{
	Hidden: true,

	Use:   "resume",
	Short: "Resumes the nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    resume,
}

// resume runs 'vagrant resume'
func resume(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	fmt.Printf(stylish.Bullet("Resuming nanobox..."))
	if err := vagrant.Resume(); err != nil {
		vagrant.Fatal("[commands/resume] vagrant.Resume() failed - ", err.Error())
	}
}
