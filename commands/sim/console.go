package sim

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
	processor.DefaultControl.Env = "sim"

	if len(args) == 0 {
		fmt.Println("you need to provide a container to console into")
		
		return
	}

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultControl.Meta["container"] = args[0]
	print.OutputCommandErr(processor.Run("sim_console", processor.DefaultControl))
}
