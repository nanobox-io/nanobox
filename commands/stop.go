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
		Short: "Stop the Nanobox virtual machine.",
		Long: `
Stops the Nanobox virtual machine as well as any running
dev and sim platforms.
		`,
		PreRun: validate.Requires("provider"),
		Run:    stopFn,
	}
)

// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Stop{}.Run())
}
