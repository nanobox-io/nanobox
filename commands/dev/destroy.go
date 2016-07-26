package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

// DestroyCmd ...
var DestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys the docker machines associated with your dev app.",
	Long: ``,
	PreRun: validate.Requires("provider", "provider_up"),
	Run:    destroyFn,
}

// destroyFn ...
func destroyFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("dev_destroy", processor.DefaultControl))
}
