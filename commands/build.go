package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	cmdutil "github.com/nanobox-io/nanobox/validate/commands"
)

var (

	// BuildCmd ...
	BuildCmd = &cobra.Command{
		Use:   "build",
		Short: "Generates a deployable build package.",
		Long: `
Generates a deployable build package that can be
deployed into a dev or production app.
		`,
		PreRun: cmdutil.Validate("provider"),
		Run:    buildFn,
	}
)

// buildFn ...
func buildFn(ccmd *cobra.Command, args []string) {

	//
	if err := processor.Run("build", processor.DefaultConfig); err != nil {

	}
}
