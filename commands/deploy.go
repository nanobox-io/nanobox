package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// DeployCmd ...
	DeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys your generated build package to a production app.",
		Long:  ``,
		Run:   deployFn,
	}

	// deployCmdFlags ...
	deployCmdFlags = struct {
		message string
	}{}
)

//
func init() {
	DeployCmd.Flags().StringVarP(&deployCmdFlags.message, "message", "m", "", "message to accompany this command")
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())

	deploy := processor.Deploy{
		Env:     env,
		App:     "default",
		Message: deployCmdFlags.message,
	}

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will return with instructions
	switch {

	// if one argument is passed we'll assume it's the name of the app to deploy to
	case len(args) == 1:
		deploy.App = args[0]

	// if more than one argument is passed we'll let the user know they are using
	// the command wrong
	case len(args) > 1:
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
name of the app you wish to deploy to:

ex: nanobox deploy <name>

`, len(args))

		return
	}

	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(deploy.Run())
}
