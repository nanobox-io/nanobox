package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// DestroyCmd ...
	DestroyCmd = &cobra.Command{
		Use:    "destroy",
		Short:  "Destroys the Nanobox virtual machine.",
		Long:   ``,
		PreRun: validate.Requires("provider"),
		Run:    destroyFn,
	}
)

// destroyFn ...
func destroyFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Destroy())
}
