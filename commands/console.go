package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// ConsoleCmd ...
	ConsoleCmd = &cobra.Command{
		Use:   "console",
		Short: "Opens an interactive console inside a production component.",
		Long:  ``,
		Run:   consoleFn,
	}

	// consoleCmdFlags ...
	consoleCmdFlags = struct {
		app 			string
		endpoint 	string
	}{}
)

//
func init() {
	ConsoleCmd.Flags().StringVarP(&consoleCmdFlags.app, "app", "a", "", "app name or alias")
	ConsoleCmd.Flags().StringVarP(&consoleCmdFlags.app, "endpoint", "e", "", "api endpoint")
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will os.Exit(1) with an error message
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
name of the component you wish to console into:

ex: nanobox console [-a appname] <container>

`, len(args))
		return
	}
	
	env, _ := models.FindEnvByID(config.EnvID())
	
	consoleConfig := processors.ConsoleConfig{
		App: consoleCmdFlags.app,
		Host: args[0],
		Endpoint: consoleCmdFlags.endpoint,
	}

	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(processors.Console(env, consoleConfig))
}
