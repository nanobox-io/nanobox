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

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
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
	if len(args) == 0 {
		args = append(args, util.Prompt("Please specify a command you wish to exec: "))
	}

	//
	v := url.Values{}
	v.Add("forward", fTunnel)
	v.Add("cmd", strings.Join(args[0:], " "))

	// fetch services to see if the command is trying to run on a specific container
	var services []Service
	if err := api.DoRawRequest(&services, "GET", fmt.Sprintf("http://%s/services", config.ServerURI), nil, nil); err != nil {
		fmt.Printf(stylish.Error("failed to get tunnel information", fmt.Sprintf("nanobox was unable to get tunnel information, and failed with the following error: %v", err)))
	}

	//
	found := false

	// range over the services to find a potential match for args[0]
	for _, service := range services {
		if args[0] == service.Name {
			found = true
		}
	}

	// if a container is found that matches args[0] then set that as a qparam, and
	// set the cmd equal to the remaining args
	if found {
		v.Add("container", args[0])
		v.Set("cmd", strings.Join(args[1:], " "))
	}

	// add a check here to regex the fTunnel to make sure the format makes sense

	docker := util.Docker{Params: v.Encode()}
	docker.Run()
}
