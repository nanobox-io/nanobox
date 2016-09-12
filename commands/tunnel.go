package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// TunnelCmd ...
	TunnelCmd = &cobra.Command{
		Use:   "tunnel",
		Short: "Creates a secure tunnel between your local machine & a production component.",
		Long: `
Creates a secure tunnel between your local machine & a
production component. The tunnel allows you to manage
production data using your local client of choice.
		`,
		PreRun: steps.Run("login"),
		Run:    tunnelFn,
	}

	// tunnelCmdFlags ...
	tunnelCmdFlags = struct {
		app  			string
		port 			string
		endpoint 	string
	}{}
)

//
func init() {
	TunnelCmd.Flags().StringVarP(&tunnelCmdFlags.app, "app", "a", "", "name or alias of a production app")
	TunnelCmd.Flags().StringVarP(&tunnelCmdFlags.port, "port", "p", "", "local port to start listening on")
	TunnelCmd.Flags().StringVarP(&tunnelCmdFlags.endpoint, "endpoint", "e", "", "api endpoint")
}

// tunnelFn ...
func tunnelFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will return with instructions
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
name of the container you would like to tunnel into:

ex: nanobox tunnel <container>

`, len(args))

		return
	}

	env, _ := models.FindEnvByID(config.EnvID())

	// set the meta arguments to be used in the processor and run the processor
	tunnelConfig := processors.TunnelConfig{
		App:       	tunnelCmdFlags.app,
		Port:      	tunnelCmdFlags.port,
		Container: 	args[0],
		Endpoint:		tunnelCmdFlags.endpoint,
	}

	// if no app id is given use 'default'
	if tunnelConfig.App == "" {
		tunnelConfig.App = "default"
	}

	display.CommandErr(processors.Tunnel(env, tunnelConfig))
}
