package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// CompileCmd ...
	CompileCmd = &cobra.Command{
		Use:   "compile",
		Short: "compile the application.",
		Long: `
Compiles the application source that can be
deployed into dev, sim, or production environments.
		`,
		PreRun: steps.Run("start", "build"),
		Run:    compileFn,
	}

)

func init() {
	steps.Build("compile", compileComplete, compileFn)
}


// compileFn ...
func compileFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	display.CommandErr(processors.Compile(env))
}

func compileComplete() bool {
	env, _ := models.FindEnvByID(config.EnvID())
	return env.Compiled
}
