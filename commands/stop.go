package commands

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
		Short: "Stop the virtual machine",
		Long: `
Stop the virtual machine.
		`,
		PreRun: validate.Requires("provider"),
		Run:    stopFn,
	}
)

// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("stop", processor.DefaultControl))
}
