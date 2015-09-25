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
	"github.com/pagodabox/nanobox-golang-stylish"
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

	// vm is running; do nothing
	case "running":
		fmt.Printf(stylish.Bullet("Nanobox VM already running"))
		break

	// vm is suspended; resume it
	case "saved":
		nanoResume(nil, args)

	// vm has not been created; create it
	case "not created":
		nanoCreate(nil, args)

	// vm is in some other state; reload just incase
	default:
		nanoReload(nil, args)
	}

	// if the vm is 'new' update images
	if util.GetVMStatus() == "not created" || !util.AppDeployed() {
		imagesUpdate(nil, args)
	}

	fRun = true
	nanoDeploy(nil, args)

	fmt.Printf("[√] App successfully built")

	var wg sync.WaitGroup

	// create a channel that listens for user interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// listen for any one event on the channel, doesn't matter what, then decrement
	// the waitgroup (blocking)
	go func() {
		<-sigChan
		close(sigChan)
		wg.Done()
	}()

	wg.Add(1)

	go nanoWatch(nil, args)

	fmt.Printf(`
----------------------------------
DEV URL : %v
----------------------------------
`, config.Nanofile.Domain)

	// set logs to streaming
	fStream = true
	go nanoLog(nil, args)

	wg.Wait()

	//
	nanoDown(nil, args)
}
