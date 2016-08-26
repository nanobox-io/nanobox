package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
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
	// get the environments
	envs, err := models.AllEnvs()
	if err != nil {
		return fmt.Printf("TODO: write message for command clean: %s\n", err.Error())
	}

	display.CommandErr(processors.Clean(envs))
}
