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
var buildCmd = &cobra.Command{
	Hidden: true,

	Use:   "build",
	Short: "Rebuilds/compiles your app",
	Long:  ``,

	PreRun:  boot,
	Run:     build,
	PostRun: halt,
}

// build
func build(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("Building codebase..."))

	// stream build output
	go Mist.Stream([]string{"log", "deploy"}, Mist.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- Mist.Listen([]string{"job", "build"}, Mist.BuildUpdates)
	}()

	// run a build
	if err := Server.Build(""); err != nil {
		Config.Fatal("[commands/build] server.Build() failed - ", err.Error())
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
