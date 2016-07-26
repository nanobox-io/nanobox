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
		Short: "Stops your sim platform.",
		Long: `
Stops your sim platform. All data and code
will be preserved in its current state.
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
