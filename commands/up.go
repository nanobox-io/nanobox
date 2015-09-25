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
	"net/http"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "",
	Long:  ``,

	Run: nanoUp,
}

//
func init() {
	upCmd.Flags().BoolVarP(&fRebuild, "rebuild", "", false, "Rebuilds")
}

//
func nanoUp(ccmd *cobra.Command, args []string) {

	//
	switch util.GetVMStatus() {

	// vm is running; do nothing
	case "running":
		fmt.Printf(stylish.Bullet("Nanobox VM already running..."))
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

	//
	switch {

	// if the vm is 'new' deploy and update images
	case util.GetVMStatus() == "not created" || !util.AppDeployed():
		imagesUpdate(nil, args)
		nanoDeploy(nil, args)

	// if fRebuild is detected only deploy
	case fRebuild:
		nanoDeploy(nil, args)
	}

	//
	nanoConsole(nil, args)

	// assume the machine can be suspended
	suspendable := true

	//
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/suspend", config.ServerURL), nil)
	if err != nil {
		config.Fatal("[commands/up] http.NewRequest() failed", err.Error())
	}

	//
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		config.Fatal("[commands/up] http.DefaultClient.Do() failed", err.Error())
	}
	defer res.Body.Close()

	//
	switch res.StatusCode / 100 {

	// anything but 200 CANNOT suspend
	default:
		suspendable = false
		fmt.Printf("\nNote: The VM has not been suspended because there there is still a console conected.\n")

	// ok to suspend
	case 2:
		break
	}

	// suspend the machine if not active consoles are connected and the command was
	// not run in background mode
	if suspendable && !fBackground {
		nanoDown(nil, args)
	}
}
