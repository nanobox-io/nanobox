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

	"github.com/spf13/cobra"

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var tunnelCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "Displays port forward information for your app's running services",
	Long: `
Description:
  Displays port forward information for your app's running services`,

	PreRun: VMIsRunning,
	Run:    nanoTunnel,
}

//
func nanoTunnel(ccmd *cobra.Command, args []string) {

	var services []Service

	//
	fmt.Printf(stylish.Bullet("Requesting services..."))
	if err := api.DoRawRequest(&services, "GET", fmt.Sprintf("http://%s/services", config.ServerURI), nil, nil); err != nil {
		fmt.Printf(stylish.Error("failed to get services", fmt.Sprintf("nanobox was unable to get services information, and failed with the following error: %v", err)))
	}

	//
	fmt.Printf(`
Service                        | Username, Password         | Ports
--------------------------------------------------------------------------------
`)
	for _, service := range services {
		fmt.Printf("%-30s | %s, %-15s | %-20v\n", service.Name, service.Username, service.Password, service.Ports)
	}
}
