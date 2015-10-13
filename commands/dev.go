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
			if err := mist.Listen([]string{"job", "deploy"}, mist.DeployUpdates); err != nil {
				config.Fatal("[commands/nanoBuild] failed - ", err.Error())
			}
			close(done)
		}()

		// wait for a status update (blocking)
		<-done
	}

	// begin watching for file changes (non blocking)
	// go func() {
	// 	if err := notify.Watch(config.CWDir, server.NotifyRebuild); err != nil {
	// 		fmt.Printf(stylish.ErrBullet("Unable to detect file changes - %v", err.Error()))
	//
	// 		// if the error is a notify error, indicate that the vm should not be suspended
	// 		if _, ok := err.(notify.WatchError); ok {
	// 			config.VMfile.SuspendableIs(false)
	// 		}
	// 	}
	// }()

	//
	server.Exec("console", "")

	// PostRun: save
}
