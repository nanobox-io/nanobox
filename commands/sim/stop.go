package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// StopCmd ...
	StopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stops the Nanobox VM and provisions app",
		Long: `
Stops the Nanobox VM and provisions app. This is the primary command stoping
the VM and preparing the application.
		`,
		PreRun: validate.Requires("provider"),
		Run:    stopFn,
	}
)

//
// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("sim_stop", processor.DefaultControl))
}
