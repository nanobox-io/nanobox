package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// CleanCmd ...
	CleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean out any apps that no longer exist",
		Long: `
todo: write long description
`,
		PreRun: validate.Requires("provider"),
		Run:    cleanFn,
	}
)

// cleanFn ...
func cleanFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("clean", processor.DefaultControl))
}
