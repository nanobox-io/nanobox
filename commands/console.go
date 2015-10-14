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
	"os"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/util/server"
)

//
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Opens an interactive terminal from inside your app on the nanobox",
	Long:  ``,

	PreRun:  boot,
	Run:     console,
	PostRun: halt,
}

// console
func console(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	//
	switch {

	// if no args are passed run console as normal
	case len(args) == 0:
		server.Exec("console", "")

	// if 1 args is passed it's assumed to be a container to console into
	case len(args) == 1:
		server.Exec("container", "container="+args[0])

	// if more than 1 args is passed fail and exit...
	case len(args) > 1:
		fmt.Printf("Expecting 0 or 1 arguments. Run 'nanobox console -h' for help. Exiting...\n")
		os.Exit(1)
	}

	// PostRun: save
}
