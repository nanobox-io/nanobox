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
	"os"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var buildCmd = &cobra.Command{
	Hidden: true,

	Use:   "build",
	Short: "Rebuilds/compiles your project",
	Long: `
Description:
  Rebuilds/compiles your project`,

	PreRun: bootVM,
	Run:    nanoBuild,
}

// nanoBuild
func nanoBuild(ccmd *cobra.Command, args []string) {

	// PreRun: bootVM

	fmt.Printf(stylish.Bullet("Building codebase..."))

	//
	build := util.Sync{
		Model:   "build",
		Path:    fmt.Sprintf("%s/builds", config.ServerURL),
		Verbose: fVerbose,
	}

	//
	build.Run(args)

	//
	switch build.Status {

	// sync completed successfully
	case "complete":
		break

	// if a build is run w/o having first run a deploy
	case "unavailable":
		fmt.Printf(stylish.ErrBullet("Before you can run a build, you must first deploy."))
		os.Exit(0)

	// errored
	case "errored":

		// this could probably be better
		if fWatch {
			config.VMfile.ModeIs("background")
		}

		//
		nanoboxDown(nil, args)
		os.Exit(1)
	}
}
