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
	"net/url"
	"os"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Opens an interactive terminal from inside your app on the nanobox VM",
	Long: `
Description:
  Opens an interactive terminal from inside your app on the nanobox VM`,

	PreRun: VMIsRunning,
	Run:    nanoConsole,
}

//
func init() {
	consoleCmd.Flags().StringVarP(&fTunnel, "tunnel", "t", "", "Creates port forwards for all comma delimeted port:port combos")
}

// nanoConsole
func nanoConsole(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Opening a nanobox console..."))

	//
	v := url.Values{}
	v.Add("forward", fTunnel)

	switch {

	// if no args are passed run console as normal
	case len(args) == 0:
		break

		// if 1 args is passed it's assumed to be a container to console directly into
		// since console doesn't take any additional arguments (like exec)
	case len(args) == 1:
		v.Add("container", args[0])

		// if more than 1 args is passed fail and exit...
	case len(args) > 1:
		fmt.Printf("Expecting 0 or 1 arguments, received %v. Exiting...\n", len(args))
		os.Exit(1)
	}

	// add a check here to regex the fTunnel to make sure the format makes sense

	docker := &util.Docker{Params: v.Encode()}
	docker.Run()
}
