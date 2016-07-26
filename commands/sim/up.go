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
		Short: "Ups the Nanobox VM and provisions your sim app.",
		Long: `
		Ups the Nanobox VM and provisions your sim app. This is the primary command uping
		the VM and preparing a sim application. It's a shortcut for 'nanobox start',
		'nanobox build', 'nanobox sim start', 'nanobox sim deploy'.
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
