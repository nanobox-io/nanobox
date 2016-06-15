package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	cmdutil "github.com/nanobox-io/nanobox/validate/commands"
)

var (

	// ResetCmd ...
	ResetCmd = &cobra.Command{
		Use:    "reset",
		Short:  "Resets the dev VM registry.",
		Long:   ``,
		PreRun: cmdutil.Validate("provider"),
		Run:    resetFn,
	}
)

// resetFn ...
func resetFn(ccmd *cobra.Command, args []string) {
	// TODO: Take an extra arguement and decide what we want to reset

	//
	if err := processor.Run("dev_reset", processor.DefaultConfig); err != nil {

	}
}
