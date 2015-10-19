// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

//
package terminal

import (
	"fmt"

	"github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
)

//
func PrintNanoboxHeader(kind string) {
	switch kind {

	//
	case "command":
		fmt.Printf(stylish.Bullet("Executing command in nanobox..."))

		//
	case "develop", "container":
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

		if kind == "develop" {
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

// GetTTYSize
func GetTTYSize(fd uintptr) (int, int) {

	ws, err := term.GetWinsize(fd)
	if err != nil {
		config.Fatal("[util/server/exec] term.GetWinsize() failed", err.Error())
	}

	//
	return int(ws.Width), int(ws.Height)
}
