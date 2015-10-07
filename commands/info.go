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

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/server"
	"github.com/nanobox-io/nanobox-cli/util/vagrant"
	// "github.com/nanobox-io/nanobox-golang-stylish"
)

//
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about the nanobox and your app",
	Long:  ``,

	Run: info,
}

// info runs 'vagrant status'
func info(ccmd *cobra.Command, args []string) {

	status := vagrant.Status()

	fmt.Printf(`
Local Domain   : %s
App Name       : %s
Nanobox state  : %s
Nanobox Files  : %s

`, config.Nanofile.Domain, config.Nanofile.Name, status, config.AppDir)

	if status != "running" {
		return
	}

	//
	var services []server.Service
	res, err := server.Get("/services", &services)
	if err != nil {
		config.Fatal("[commands/nanoboxInfo] failed - ", err.Error())
	}
	defer res.Body.Close()

	//
	if len(services) >= 1 {
		info := "///////// SERVICES /////////\n\n"

		//
		for _, service := range services {
			info += fmt.Sprintf("%s :\n", service.UID)

			if service.Name != "" {
				info += fmt.Sprintf("   name : %s\n", service.Name)
			}

			info += fmt.Sprintf("   host : %s\n   ports : %v\n", config.Nanofile.Domain, service.Ports)

			//
			if service.Username != "" {
				info += fmt.Sprintf("   username : %s\n", service.Username)
			}

			//
			if service.Password != "" {
				info += fmt.Sprintf("   password : %s\n", service.Password)
			}
		}

		//
		fmt.Printf("%s\n\n", info)
	}
}
