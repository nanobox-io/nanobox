package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// BuildCmd ...
	BuildCmd = &cobra.Command{
		Use:   "build-runtime",
		Short: "Builds a deployable runtime.",
		Long: `
Generates a deployable runtime that can be
deployed into dev, sim, or production environments.
		`,
		PreRun: steps.Run("start"),
		Run:    buildFn,
		Aliases: []string{"build"},
	}
)

func init() {
	steps.Build("build-runtime", buildComplete, buildFn)
}

// buildFn ...
func buildFn(ccmd *cobra.Command, args []string) {

	env, _ := models.FindEnvByID(config.EnvID())
	display.CommandErr(processors.Build(env))
}

func buildComplete() bool {
	env, _ := models.FindEnvByID(config.EnvID())
	box := boxfile.NewFromPath(config.Boxfile())

	return env.UserBoxfile != "" && env.UserBoxfile == box.String()
}
