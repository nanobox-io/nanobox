package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// BuildCmd ...
	BuildCmd = &cobra.Command{
		Use:   "build",
		Short: "Generates a deployable build package.",
		Long: `
Generates a deployable build package that can be
deployed into dev, sim, or production platforms.
		`,
		PreRun: validate.Requires("provider"),
		Run:    buildFn,
	}
)

// buildFn ...
func buildFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("build", processor.DefaultControl))
}
