package dev

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	cmdutil "github.com/nanobox-io/nanobox/validate/commands"
)

// ConsoleCmd ...
var ConsoleCmd = &cobra.Command{
	Use:    "console",
	Short:  "Opens an interactive console inside your Nanobox VM.",
	Long:   ``,
	PreRun: cmdutil.Validate("provider"),
	Run:    consoleFn,
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("I need a container to run in")
		return
	}
	processor.DefaultConfig.Meta["name"] = args[0]

	//
	if err := processor.Run("dev_console", processor.DefaultConfig); err != nil {

	}
}
