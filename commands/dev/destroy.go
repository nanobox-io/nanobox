package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	cmdutil "github.com/nanobox-io/nanobox/validate/commands"
)

// DestroyCmd ...
var DestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys the docker machines associated with the current app.",
	Long: `
Destroys the docker machines associated with the current app.
If no other apps are running, it will destroy the Nanobox VM.
		`,
	PreRun: cmdutil.Validate("provider"),
	Run:    destroyFn,
}

// destroyFn ...
func destroyFn(ccmd *cobra.Command, args []string) {

	//
	if err := processor.Run("dev_destroy", processor.DefaultConfig); err != nil {

	}
}
