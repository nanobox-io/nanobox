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

	buildNoCompile bool
)

func init() {
	BuildCmd.PersistentFlags().BoolVarP(&buildNoCompile, "no-compile", "", false, "dont compile the build")
}

// buildFn ...
func buildFn(ccmd *cobra.Command, args []string) {
	if buildNoCompile {
		processor.DefaultControl.Meta["no-compile"] = "true"
	}
	print.OutputCommandErr(processor.Run("build", processor.DefaultControl))
}
