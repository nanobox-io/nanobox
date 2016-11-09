package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"

	// added because we need its steps
	_ "github.com/nanobox-io/nanobox/commands/sim"
)

var (

	// DeployCmd ...
	DeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your application to a live remote or a dry-run environment.",
		Long:  ``,
		PreRun: func(ccmd *cobra.Command, args []string) {
			registry.Set("skip-compile", deployCmdFlags.skipCompile)
			steps.Run("configure", "start", "build-runtime", "compile-app")(ccmd, args)
		},
		Run: deployFn,
	}

	// deployCmdFlags ...
	deployCmdFlags = struct {
		skipCompile bool
		message     string
		force       bool
	}{}
)

//
func init() {
	DeployCmd.Flags().BoolVarP(&deployCmdFlags.skipCompile, "skip-compile", "", false, "skip compiling the app")
	DeployCmd.Flags().BoolVarP(&deployCmdFlags.force, "force", "", false, "force the deploy even if you have used this build on a previous deploy")
	DeployCmd.Flags().StringVarP(&deployCmdFlags.message, "message", "m", "", "message to accompany this command")
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {
	envModel, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(envModel, args)

	switch location {
	case "local":
		switch name {
		case "dev":
			fmt.Println("deploying is not necessary in this context, 'nanobox run' instead")
			return
		case "sim":
			steps.Run("sim start")(ccmd, args)
			appModel, _ := models.FindAppBySlug(envModel.ID, "sim")
			display.CommandErr(app.Deploy(envModel, appModel))
			steps.Run("sim stop")(ccmd, args)
		}
	case "production":
		steps.Run("login")(ccmd, args)
		deployConfig := processors.DeployConfig{
			App:      name,
			Message:  deployCmdFlags.message,
			Force:    deployCmdFlags.force,
		}

		// set the meta arguments to be used in the processor and run the processor
		display.CommandErr(processors.Deploy(envModel, deployConfig))
	}
}
