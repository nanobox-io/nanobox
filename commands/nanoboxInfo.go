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

	api "github.com/nanobox-io/nanobox-api-client"
	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	// "github.com/nanobox-io/nanobox-golang-stylish"
)

//
var nanoboxInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about the nanobox and your app",
	Long:  ``,

	Run: nanoboxInfo,
}

// nanoboxInfo runs 'vagrant status'
func nanoboxInfo(ccmd *cobra.Command, args []string) {

	status := util.VagrantStatus()

	fmt.Printf(`
Local Domain   : %s
App Name       : %s
Nanobox state  : %s
Nanobox Files  : %s

`, config.Nanofile.Domain, config.Nanofile.Name, status, config.AppDir)

	if status != "running" {
		return
	}

	var services []Service

	//
	if err := api.DoRawRequest(&services, "GET", fmt.Sprintf("%s/services", config.ServerURL), nil, nil); err != nil {
		panic(err)
	}

	//
	if len(services) >= 1 {
		info := "///////// SERVICES /////////\n"

		//
		for _, service := range services {
			info += fmt.Sprintf(`
  %s :
    host : %s
    ports : %v
			`, service.Name, config.Nanofile.Domain, service.Ports)

			//
			if service.Username != "" {
				info += fmt.Sprintf("username : %s", service.Username)
			}

			//
			if service.Password != "" {
				info += fmt.Sprintf("password : %s", service.Password)
			}
		}

		//
		fmt.Println(info)
	}
}
