package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// StartCmd ...
	StartCmd = &cobra.Command{
		Use:    "start",
		Short:  "Starts the Nanobox virtual machine.",
		Long:   ``,
		PreRun: validate.Requires("provider"),
		Run:    startFn,
	}
)

// startFn ...
func startFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Start{}.Run())
}
