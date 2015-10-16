// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

//
package server

import (
	"fmt"

	terminal "github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

// IsContainerExec
func IsContainerExec(args []string) (found bool) {

	// fetch services to see if the command is trying to run on a specific container
	var services []Service
	res, err := Get("/services", &services)
	if err != nil {
		Fatal("[util/server/exec] Get() failed - ", err.Error())
	}
	defer res.Body.Close()

	//
	for _, service := range services {

		// range over the services to find a potential match for args[0] or make an
		// exception for 'build1' since that wont show up on the list.
		if args[0] == service.Name || args[0] == "build1" {
			found = true
		}
	}

	return
}

//
func header(kind string) {
	switch kind {

	//
	case "command":
		fmt.Printf(stylish.Bullet("Executing command in nanobox..."))

		//
	case "console", "container":
		fmt.Printf(`+> Opening a nanobox console:


                                     **
                                  ********
                               ***************
                            *********************
                              *****************
                            ::    *********    ::
                               ::    ***    ::
                             ++   :::   :::   ++
                                ++   :::   ++
                                   ++   ++
                                      +

                      _  _ ____ _  _ ____ ___  ____ _  _
                      |\ | |__| |\ | |  | |__) |  |  \/
                      | \| |  | | \| |__| |__) |__| _/\_
`)

		if kind == "console" {
			fmt.Printf(`
--------------------------------------------------------------------------------
+ You are in a virtual machine (vm)
+ Your local source code has been mounted into the vm, and changes in either
the vm or local will be mirrored.
+ If you run a server, access it at >> %s
--------------------------------------------------------------------------------
`, config.Nanofile.Domain)
		}
	}
}

// getTTYSize
func getTTYSize(fd uintptr) (int, int) {

	ws, err := terminal.GetWinsize(fd)
	if err != nil {
		config.Fatal("[util/server/exec] terminal.GetWinsize() failed", err.Error())
	}

	//
	if ws == nil {
		return 0, 0
	}

	//
	return int(ws.Width), int(ws.Height)
}
