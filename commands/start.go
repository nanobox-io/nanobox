package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
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
	print.OutputCommandErr(processor.Start{}.Run())
}
