package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// TunnelCmd ...
	TunnelCmd = &cobra.Command{
		Use:   "tunnel",
		Short: "Create a secure tunnel between your local machine & a live component.",
		Long: `
Creates a secure tunnel between your local machine &
a live component. The tunnel allows you to manage
live data using your local client of choice.
		`,
		PreRun: steps.Run("login"),
		Run:    tunnelFn,
	}

	// tunnelCmdFlags ...
	tunnelCmdFlags = struct {
		port     string
	}{}
)

//
func init() {
	TunnelCmd.Flags().StringVarP(&tunnelCmdFlags.port, "port", "p", "", "local port to start listening on")
}

// tunnelFn ...
func tunnelFn(ccmd *cobra.Command, args []string) {

	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args)

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


	switch location {
	case "local":
		fmt.Println("tunneling is not required for local development")
		return
	case "production":
		// set the meta arguments to be used in the processor and run the processor
		tunnelConfig := processors.TunnelConfig{
			App:       name,
			Port:      tunnelCmdFlags.port,
			Container: args[0],
		}

		display.CommandErr(processors.Tunnel(env, tunnelConfig))
	}

}
