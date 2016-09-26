package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// ImplodeCmd ...
	ImplodeCmd = &cobra.Command{
		Use:    "implode",
		Short:  "remove all nanobox created files and vms",
		Long:   ``,
		PreRun: steps.Run("start"),
		Run:    implodeFn,
	}
)

// implodeFn ...
func implodeFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Implode())
}
