// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"os"
	"os/signal"
	"sync"

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
	// issue a deploy --run, watch for file changes, and display streaming logs
	case fRun:
		var wg sync.WaitGroup

		// create a channel that listens for user interrupts
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		// listen for anything on the channel (doesn't matter what) and exit
		go func() {
			for _ = range sigChan {
				wg.Done()
				// os.Exit(0)
			}
		}()

		nanoCreate(nil, args)
		imagesUpdate(nil, args)
		nanoDeploy(nil, args)

		wg.Add(1)

		// set logs to streaming
		fStream = true

		go nanoLog(nil, args)
		go nanoWatch(nil, args)

		wg.Wait()

	}
}
