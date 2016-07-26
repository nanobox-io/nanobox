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
		Short: "Stops your dev platform.",
		Long: `
Stops your dev platform. All data will be preserved in its current state.
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
