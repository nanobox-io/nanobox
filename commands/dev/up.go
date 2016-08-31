package dev

import (
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
	"github.com/spf13/cobra"
)

var (

	// UpCmd ...
	UpCmd = &cobra.Command{
		Use:   "up",
		Short: "Ups the Nanobox VM and provisions your dev app.",
		Long: `
Ups the Nanobox VM and provisions your dev app. This is the primary command uping
the VM and preparing a dev application. It's a shortcut for 'nanobox start',
'nanobox build', 'nanobox dev start', 'nanobox dev deploy'.
		`,
		PreRun: validate.Requires("provider"),
		Run:    upFn,
	}
)

// upFn ...
func upFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(dev.Up())
}
