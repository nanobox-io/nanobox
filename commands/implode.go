package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// ImplodeCmd ...
	ImplodeCmd = &cobra.Command{
		Use:   "implode",
		Short: "Remove all Nanobox-created containers, files, & data.",
		Long: `
Removes the Nanobox container, all projects, filesystem mounts,
& local data. All that will remain is nanobox binaries.
		`,
		Run: implodeFn,
	}
)

// implodeFn ...
func implodeFn(ccmd *cobra.Command, args []string) {
	registry.Set("full-implode", true)
	display.CommandErr(processors.Implode())
}
