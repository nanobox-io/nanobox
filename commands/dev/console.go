//
package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/server"
)

//
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Opens an interactive terminal from inside your app on nanobox",
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
		fmt.Printf(stylish.ErrBullet("Unable to console. Please provide a service to connect to.\n"))

	// if 1 args is passed it's assumed to be a container to console into
	case len(args) == 1:
		if err := server.Console("container=" + args[0]); err != nil {
			server.Error("[commands/console] Server.Console failed", err.Error())
		}
	}

	// PostRun: halt
}
