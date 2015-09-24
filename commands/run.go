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

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
)

//
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  ``,

	Run: nanoRun,
}

//
func init() {
	runCmd.Flags().BoolVarP(&fForce, "reset-cache", "", false, "resets stuff")
}

//
func nanoRun(ccmd *cobra.Command, args []string) {

	//
	switch util.GetVMStatus() {
	case "not created":
		nanoCreate(nil, args)
	case "saved":
		nanoResume(nil, args)
	default:
		nanoReload(nil, args)
	}

	fRun = true
	nanoDeploy(nil, args)

	fmt.Printf("[√] App successfully built")

	var wg sync.WaitGroup

	// create a channel that listens for user interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// listen for anything on the channel (doesn't matter what) and exit
	go func() {
		for _ = range sigChan {
			wg.Done()
		}
	}()

	wg.Add(1)

	// set logs to streaming
	fStream = true

	go nanoWatch(nil, args)

	fmt.Printf(`
----------------------------------
DEV URL : %v
----------------------------------

`, config.Nanofile.Domain)

	go nanoLog(nil, args)

	wg.Wait()

	// suspend the machine if not debug mode
	if !fDebug {
		nanoDown(nil, args)
	}
}
