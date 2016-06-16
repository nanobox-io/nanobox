package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
)

var (

	// DeployCmd ...
	DeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys your generated build package to a production app.",
		Long:  ``,
		Run:   deployFn,
	}
)

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will os.Exit(1) with an error message
	switch {

	// if no arguments are passed we'll deploy to the "default" app...
	case len(args) == 0:
		processor.DefaultConfig.Meta["app_name"] = "default"

	// if one argument is passed we'll assume it's the name of the app to deploy to
	case len(args) == 1:
		processor.DefaultConfig.Meta["app_name"] = args[0]

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
	print.OutputCommandErr(processor.Run("deploy", processor.DefaultConfig))
}
