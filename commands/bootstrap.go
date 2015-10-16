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
var bootstrapCmd = &cobra.Command{
	Hidden: true,

	Use:   "bootstrap",
	Short: "Runs an engine's bootstrap script - downloads code & launches a nanobox",
	Long:  ``,

	PreRun:  boot,
	Run:     bootstrap,
	PostRun: halt,
}

// bootstrap
func bootstrap(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("Bootstrapping codebase..."))

	// stream bootstrap output
	go Mist.Stream([]string{"log", "deploy"}, Mist.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- Mist.Listen([]string{"job", "bootstrap"}, Mist.BootstrapUpdates)
	}()

	// run a bootstrap
	if err := Server.Bootstrap(""); err != nil {
		Config.Fatal("[commands/bootstrap] server.Bootstrap() failed - ", err.Error())
	}

	// wait for a status update (blocking)
	err := <-errch

	//
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// PostRun: halt
}
