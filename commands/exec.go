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
	"strings"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Runs a command from inside your app on the nanobox VM",
	Long: `
Description:
  Runs a command from inside your app on the nanobox VM`,

	Run: nanoExec,
}

//
func init() {
	execCmd.Flags().StringVarP(&fTunnel, "tunnel", "t", "", "Creates port forwards for all comma delimeted port:port combos")
}

// nanoExec
func nanoExec(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Opening a nanobox console..."))

	//
	if len(args) <= 0 {
		args = append(args, ui.Prompt("Please specify a command you wish to exec: "))
	}

	// add a check here to regex the fTunnel to make sure the format makes sense

	//
	v := url.Values{}
	v.Add("forward", fTunnel)
	v.Add("cmd", strings.Join(args, " "))

	docker := utils.Docker{Params: v.Encode()}
	docker.Run()
}
