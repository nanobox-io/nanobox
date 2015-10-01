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

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	// "github.com/nanobox-io/nanobox-golang-stylish"
)

//
var nanoboxConsoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Opens an interactive terminal from inside your app on the nanobox",
	Long:  ``,

	PreRun:  bootVM,
	Run:     nanoboxConsole,
	PostRun: saveVM,
}

// nanoboxConsole
func nanoboxConsole(ccmd *cobra.Command, args []string) {

	// PreRun: bootVM

	msg := `+> Opening a nanobox console:


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
`

	//
	v := url.Values{}

	//
	switch {

	// if no args are passed run console as normal
	case len(args) == 0:
		msg += fmt.Sprintf(`
------------------------------------------------------------------
+ You are in a virtual machine (vm)
+ Your local source code has been mounted into the vm, and changes
 in either the vm or local will be mirrored.
+ If you run a server, access it at >> %s
------------------------------------------------------------------
	`, config.Nanofile.Domain)
		break

		// if 1 args is passed it's assumed to be a container to console directly into
		// since console doesn't take any additional arguments (like exec)
	case len(args) == 1:
		v.Add("container", args[0])

		// if more than 1 args is passed fail and exit...
	case len(args) > 1:
		fmt.Printf("Expecting 0 or 1 arguments. Exiting...\n")
		os.Exit(1)
	}

	fmt.Println(msg)

	//
	server := &util.Server{Params: v.Encode()}
	server.Run()

	// PostRun: saveVM
}
