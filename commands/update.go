package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// UpdateCmd ...
	UpdateCmd = &cobra.Command{
		Use:    "update",
		Short:  "Updates docker images and checks to see if the nanobox binary needs an update.",
		Long:   ``,
		PreRun: steps.Run("start"),
		Run:    updateFn,
	}
)

// updateFn ...
func updateFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Update())
}
