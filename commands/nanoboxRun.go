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
	"os/signal"
	"sync"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	// "github.com/nanobox-io/nanobox-cli/util"
	// "github.com/nanobox-io/nanobox-golang-stylish"
)

//
var nanoboxRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts a nanobox, provisions the app, & runs the app's exec",
	Long:  ``,

	PreRun:  bootVM,
	Run:     nanoboxRun,
	PostRun: saveVM,
}

//
func init() {
	nanoboxRunCmd.Flags().BoolVarP(&fForce, "reset-cache", "", false, "resets stuff")
}

// nanoboxRun
func nanoboxRun(ccmd *cobra.Command, args []string) {

	// PreRun: bootVM

	fRun = true
	nanoDeploy(nil, args)

	var wg sync.WaitGroup

	// create a channel that listens for user interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// listen for any one event on the channel, doesn't matter what, then decrement
	// the waitgroup (blocking)
	go func() {
		<-sigChan
		wg.Done()
		close(sigChan)
	}()

	wg.Add(1)

	go nanoWatch(nil, args)

	// set logs to streaming
	fStream = true
	go nanoLog(nil, args)

	fmt.Printf(`
--------------------------------------------------------------------------------
[√] APP SUCCESSFULLY BUILT   ///   DEV URL : %v
--------------------------------------------------------------------------------

++> STREAMING LOGS (ctrl-c to exit) >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
`, config.Nanofile.Domain)

	wg.Wait()

	// PostRun: saveVM
}
