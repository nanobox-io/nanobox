package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
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
		PreRun: steps.Run("start"),
		Run:    buildFn,
	}

	buildSkipCompile bool
)

func init() {
	BuildCmd.PersistentFlags().BoolVarP(&buildSkipCompile, "skip-compile", "", false, "dont compile the build")

	steps.Build("build", buildComplete, buildFn)
}

// buildFn ...
func buildFn(ccmd *cobra.Command, args []string) {

	if buildSkipCompile {
		registry.Set("skip-compile", true)
	}

	env, _ := models.FindEnvByID(config.EnvID())
	display.CommandErr(processors.Build(env))
}

func buildComplete() bool {
	env, _ := models.FindEnvByID(config.EnvID())
	box := boxfile.NewFromPath(config.Boxfile())

	return env.UserBoxfile != "" && env.UserBoxfile == box.String()
}
