package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// CleanCmd ...
	CleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean out any environemnts that no longer exist",
		Long: `
todo: write long description
`,
		PreRun: validate.Requires("provider"),
		Run:    cleanFn,
	}
)

// cleanFn ...
func cleanFn(ccmd *cobra.Command, args []string) {
	clean := processor.Clean{}
	display.CommandErr(clean.Run())
}
