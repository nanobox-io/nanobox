// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"github.com/spf13/cobra"
)

//
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "",
	Long: `
Description:
  Runs 'nanobox create' and then 'nanobox deploy'`,

	Run: nanoUp,
}

//
func init() {
	upCmd.Flags().BoolVarP(&fRun, "run", "", false, "Watches your app for file changes")
}

//
func nanoUp(ccmd *cobra.Command, args []string) {

	switch {

	// by default, create the environment, update all images, issue a deploy and
	// drop the user into a console
	default:
		nanoCreate(nil, args)
		imagesUpdate(nil, args)
		nanoDeploy(nil, args)
		nanoConsole(nil, args)

	// if the --run flag is found, create the environment, update docker images,
	// issue a deploy --run, and watch for file changes, and show logs
	case fRun:
		nanoCreate(nil, args)
		imagesUpdate(nil, args)
		nanoDeploy(nil, []string{"--run"})
		go nanoWatch(nil, args)
		go nanoLog(nil, []string{"--stream"})
	}
}
