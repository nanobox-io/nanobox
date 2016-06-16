package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// ResetCmd ...
	ResetCmd = &cobra.Command{
		Use:    "reset",
		Short:  "Resets the dev VM registry.",
		Long:   ``,
		PreRun: validate.Requires("provider"),
		Run:    resetFn,
	}
)

// TODO: Take an extra arguement and decide what we want to reset
// resetFn ...
func resetFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("dev_reset", processor.DefaultConfig))
}
