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

	"github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Opens an interactive terminal from inside your app on the nanobox VM",
	Long: `
Description:
  Opens an interactive terminal from inside your app on the nanobox VM`,

	Run: nanoConsole,
}

//
func init() {
	consoleCmd.Flags().StringVarP(&fTunnel, "tunnel", "t", "", "Creates port forwards for all comma delimeted port:port combos")
}

// nanoConsole
func nanoConsole(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Opening a nanobox console..."))

	//
	if len(args) > 0 {
		fmt.Println("Attempting to run 'nanobox console' with a command. Use 'nanobox exec'")
		os.Exit(0)
	}

	// add a check here to regex the fTunnel to make sure the format makes sense

	//
	v := url.Values{}
	v.Add("forward", fTunnel)

	docker := utils.Docker{Params: v.Encode()}
	docker.Run()
}
