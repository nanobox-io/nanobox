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
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "",
	Long:  ``,

	Run: nanoInfo,
}

// nanoInfo runs 'vagrant status'
func nanoInfo(ccmd *cobra.Command, args []string) {

	status := util.GetVMStatus()

	fmt.Printf(`
Local Domain   : %s
App Name       : %s
VM state       : %s
Nanobox Files  : %s
`, config.Nanofile.Domain, config.Nanofile.Name, status, config.AppDir)

	if status != "running" {
		return
	}

	var services []Service

	//
	if err := api.DoRawRequest(&services, "GET", fmt.Sprintf("http://%s/services", config.ServerURI), nil, nil); err != nil {
		fmt.Printf(stylish.Error("failed to get services", fmt.Sprintf("nanobox was unable to get services information, and failed with the following error: %v", err)))
	}

	//
	if len(services) >= 1 {
		fmt.Printf("///////// SERVICES /////////\n")

		//
		for _, service := range services {
			info := fmt.Sprintf(`
	%s :
	  name : %s
		host : %s
		ports : %v
			`, service.UID, service.Name, config.Nanofile.Domain, service.Ports)

			//
			if service.Username != "" {
				info += fmt.Sprintf("username : %s", service.Username)
			}

			//
			if service.Password != "" {
				info += fmt.Sprintf("password : %s", service.Password)
			}
		}
	}
}

// ///////// ENV VARS /////////
//
// somevar : "nothing"
// var1    : "hello"
// var2    : "world"
