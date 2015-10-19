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

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/server"
	"github.com/nanobox-io/nanobox/util/server/mist"
	"github.com/nanobox-io/nanobox/util/vagrant"
)

var (

	//
	devCmd = &cobra.Command{
		Use:   "dev",
		Short: "Starts the nanobox, provisions app, & opens an interactive terminal",
		Long:  ``,

		PreRun:  boot,
		Run:     dev,
		PostRun: halt,
	}

	//
	rebuild bool // force a deploy
	nobuild bool // force skip a deploy
)

//
func init() {
	devCmd.Flags().BoolVarP(&rebuild, "rebuild", "", false, "Force a rebuild")
	devCmd.Flags().BoolVarP(&nobuild, "no-build", "", false, "Force skip a rebuild")
}

// dev
func dev(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	if !nobuild {

		// if the vm has no been created, deployed, or the rebuild flag is passed do
		// a deploy
		if vagrant.Status() == "not created" || !config.VMfile.HasDeployed() || rebuild {

			fmt.Printf(stylish.Bullet("Deploying codebase..."))

			// run a deploy
			if err := server.Deploy(""); err != nil {
				server.Fatal("[commands/dev] server.Deploy() failed - ", err.Error())
			}

			// stream log output
			go mist.Stream([]string{"log", "deploy"}, mist.PrintLogStream)

			// listen for status updates
			errch := make(chan error)
			go func() {
				errch <- mist.Listen([]string{"job", "deploy"}, mist.DeployUpdates)
			}()

			// wait for a status update (blocking)
			err := <-errch

			//
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
		}
	}

	//
	server.Exec("develop", "console", "")

	// PostRun: halt
}
