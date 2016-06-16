package dev

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

// ConsoleCmd ...
var ConsoleCmd = &cobra.Command{
	Use:    "console",
	Short:  "Opens an interactive console inside your Nanobox VM.",
	Long:   ``,
	PreRun: validate.Requires("provider"),
	Run:    consoleFn,
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will os.Exit(1) with instructions
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
name of the container you wish to console into:

ex: nanobox dev console <container>

`, len(args))
		return
	}

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultConfig.Meta["container"] = args[0]
	print.OutputCommandErr(processor.Run("dev_console", processor.DefaultConfig))
}
