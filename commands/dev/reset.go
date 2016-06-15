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

// resetFn ...
func resetFn(ccmd *cobra.Command, args []string) {
	// TODO: Take an extra arguement and decide what we want to reset

	//
	print.OutputCmdErr(processor.Run("dev_reset", processor.DefaultConfig))
}
