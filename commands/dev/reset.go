package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/display"
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
		Hidden: true,
	}
)

// TODO: Take an extra arguement and decide what we want to reset?
// resetFn ...
func resetFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(dev.Reset())
}
