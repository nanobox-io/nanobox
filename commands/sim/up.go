package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// UpCmd ...
	UpCmd = &cobra.Command{
		Use:   "up",
		Short: "Ups the Nanobox VM and provisions app",
		Long: `
Ups the Nanobox VM and provisions app.
		`,
		PreRun: validate.Requires("provider"),
		Run:    simUp,
	}
)

//
// simUp ...
func simUp(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("sim_up", processor.DefaultControl))
}
