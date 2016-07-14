package dev

import (
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
	PreRun: validate.Requires("provider", "provider_up", "dev_isup"),
	Run:    consoleFn,
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {
	processor.DefaultControl.Env = "dev"

	// if given an argument they wanted to run a console into a container
	// if no arguement is provided they wanted to run a dev console
	// and be dropped into a dev environment
	if len(args) > 0 {
		processor.DefaultControl.Meta["container"] = args[0]
		print.OutputCommandErr(processor.Run("env_console", processor.DefaultControl))
		return
	}

	// set the meta arguments to be used in the processor and run the processor
	print.OutputCommandErr(processor.Run("dev_console", processor.DefaultControl))
}
