package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// UpCmd ...
	UpCmd = &cobra.Command{
		Use:   "up",
		Short: "Ups the Nanobox VM and provisions app",
		Long: `
Ups the Nanobox VM and provisions app. This is the primary command uping
the VM and preparing the application.
		`,
		PreRun: validate.Requires("provider"),
		Run:    devUp,
	}
)

//
// devUp ...
func devUp(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("dev_up", processor.DefaultControl))
}
