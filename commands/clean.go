package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// CleanCmd ...
	CleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean out any environemnts that no longer exist",
		Long: `
todo: write long description
`,
		Run:    cleanFn,
	}
)

// cleanFn ...
func cleanFn(ccmd *cobra.Command, args []string) {
	// get the environments
	envs, err := models.AllEnvs()
	if err != nil {
		fmt.Printf("TODO: write message for command clean: %s\n", err.Error())
		return
	}

	display.CommandErr(processors.Clean(envs))
}
