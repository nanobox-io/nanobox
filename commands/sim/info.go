package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// InfoCmd ...
	InfoCmd = &cobra.Command{
		Use:    "info",
		Short:  "Displays information about the running sim app and its components.",
		Long:   ``,
		PreRun: validate.Requires("provider"),
		Run:    infoFn,
	}
)

// infoFn will run the DNS processor for adding DNS entires to the "hosts" file
func infoFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("sim_info", processor.DefaultControl))
}
