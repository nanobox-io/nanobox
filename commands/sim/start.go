package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// StartCmd ...
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts your sim platform.",
		Long: `
Starts the sim platform from its previous state. If starting for
the first time, you should also generate a build (nanobox build)
and deploy it into your sim platform (nanobox sim deploy).
		`,
		PreRun: validate.Requires("provider"),
		Run:    startFn,
	}
)

//
// startFn ...
func startFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("sim_start", processor.DefaultControl))
}
