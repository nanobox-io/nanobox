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

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Rebuilds/compiles your project",
	Long: `
Description:
  Rebuilds/compiles your project`,

	Run: nanoBuild,
}

// nanoBuild
func nanoBuild(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Building codebase..."))

	//
	build := util.Sync{
		Model:   "build",
		Path:    fmt.Sprintf("http://%v:1757/builds", config.Nanofile.IP),
		Verbose: fVerbose,
	}

	//
	build.Run(args)

	//
	switch build.Status {

	// complete
	case "complete":
		fmt.Printf(stylish.Bullet(fmt.Sprintf("Build complete... Navigate to %v.nano.dev to view your app.", config.App)))

		// if the build fails the server should handle the message. If not, this can
		// be re-enabled
	case "errored":
		// fmt.Printf(stylish.Error("Build failed", "Your build failed to well... build"))
	}
}
