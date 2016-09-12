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

	// DeployCmd ...
	DeployCmd = &cobra.Command{
		Use:    "deploy",
		Short:  "Deploys your generated build package to a production app.",
		Long:   ``,
		PreRun: steps.Run("start", "build", "login"),
		Run:    deployFn,
	}

	// deployCmdFlags ...
	deployCmdFlags = struct {
		app      string
		message  string
		force    bool
		endpoint string
	}{}
)

//
func init() {
	DeployCmd.Flags().BoolVarP(&deployCmdFlags.force, "force", "", false, "force the deploy even if you have used this build on a previous deploy")
	DeployCmd.Flags().StringVarP(&deployCmdFlags.app, "app", "a", "", "message to accompany this command")
	DeployCmd.Flags().StringVarP(&deployCmdFlags.message, "message", "m", "", "message to accompany this command")
	DeployCmd.Flags().StringVarP(&deployCmdFlags.endpoint, "endpoint", "e", "", "api endpoint")
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	// TODO: make sure the environmetn is setup

	deployConfig := processors.DeployConfig{
		App:      deployCmdFlags.app,
		Message:  deployCmdFlags.message,
		Force:    deployCmdFlags.force,
		Endpoint: deployCmdFlags.endpoint,
	}

	if deployConfig.App == "" {
		deployConfig.App = "default"
	}

	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(processors.Deploy(env, deployConfig))
}
