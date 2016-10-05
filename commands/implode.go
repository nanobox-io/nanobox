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
		Use:   "implode",
		Short: "Removes all Nanobox-created containers, files, & data",
		Long: `
Removes the Nanobox container, all projects, filesystem mounts,
& local data. All that will remain is nanobox binaries.
		`,
		PreRun: steps.Run("start"),
		Run:    implodeFn,
	}
)

// implodeFn ...
func implodeFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Implode())
}
