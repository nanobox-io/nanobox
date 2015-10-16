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

	// if no args are passed provide instruction
	case len(args) == 0:
		fmt.Printf("Please provide a service to connect to.\n")

	// if 1 args is passed it's assumed to be a container to console into
	case len(args) == 1:
		server.Exec("container", "container="+args[0])
	}

	// PostRun: halt
}
