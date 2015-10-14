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
	"github.com/nanobox-io/nanobox-golang-stylish"
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
	go mist.Stream([]string{"log", "deploy"}, mist.PrintLogStream)

	// listen for status updates
	done := make(chan struct{})
	go func() {
		if err := mist.Listen([]string{"job", "build"}, mist.BuildUpdates); err != nil {
			config.Fatal("[commands/build] failed - ", err.Error())
		}
		close(done)
	}()

	// run a build
	if err := server.Build(""); err != nil {
		config.Fatal("[commands/build] failed - ", err.Error())
	}

	// wait for a status update (blocking)
	<-done

	// PostRun: save
}
