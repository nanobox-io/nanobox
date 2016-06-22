package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// StartCmd ...
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the Nanobox VM and provisions app",
		Long: `
Starts the Nanobox VM and provisions app. This is the primary command starting
the VM and preparing the application.
		`,
		PreRun: validate.Requires("provider"),
		Run:    devStart,
	}
)

//
// devStart ...
func devStart(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("dev_start", processor.DefaultControl))
}
