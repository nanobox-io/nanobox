package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor/sim"
	"github.com/nanobox-io/nanobox/util/display"
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
		Run:    upFn,
	}
)

// upFn ...
func upFn(ccmd *cobra.Command, args []string) {
	simUp := sim.Up{}
	display.CommandErr(simUp.Run())
}
