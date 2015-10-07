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

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/server"
	"github.com/nanobox-io/nanobox-cli/util/server/mist"
	"github.com/nanobox-io/nanobox-cli/util/vagrant"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Starts the nanobox, provisions app, & opens an interactive terminal",
	Long:  ``,

	PreRun:  boot,
	Run:     dev,
	PostRun: save,
}

//
func init() {
	devCmd.Flags().BoolVarP(&fRebuild, "rebuild", "", false, "Rebuilds")
}

// dev
func dev(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	// if the vm has no been created, deployed, or the rebuild flag is passed do
	// a deploy
	if vagrant.Status() == "not created" || !config.VMfile.HasDeployed() || fRebuild {

		fmt.Printf(stylish.Bullet("Deploying codebase..."))

		// run a deploy
		if err := server.Deploy(""); err != nil {
			config.Fatal("[commands/nanoDeploy] failed - ", err.Error())
		}

		// stream log output
		go mist.Stream([]string{"log", "deploy"}, mist.PrintLogStream)

		// listen for status updates
		done := make(chan struct{})
		go func() {
			if err := mist.Listen([]string{"job", "deploy"}, mist.HandleDeployStream); err != nil {
				config.Fatal("[commands/nanoBuild] failed - ", err.Error())
			}
			close(done)
		}()

		// wait for a status update (blocking)
		<-done
	}

	//
	server.Exec("console", "")

	// PostRun: save
}
