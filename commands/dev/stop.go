package dev

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
		Short: "Stops the Nanobox app",
		Long: `
Stops the Nanobox app.
		`,
		PreRun: validate.Requires("provider"),
		Run:    devStop,
	}
)

//
// devStop ...
func devStop(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("dev_stop", processor.DefaultControl))
}
