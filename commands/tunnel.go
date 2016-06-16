package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (
	port string

	// TunnelCmd ...
	TunnelCmd = &cobra.Command{
		Use:   "tunnel",
		Short: "Creates a secure tunnel between your local machine & a production component.",
		Long: `
Creates a secure tunnel between your local machine & a
production component. The tunnel allows you to manage
production data using your local client of choice.
		`,
		PreRun: validate.Requires("provider"),
		Run:    tunnelFn,
	}
)

//
func init() {
	TunnelCmd.Flags().StringVarP(&app, "app", "a", "", "production app name or alias")
	TunnelCmd.Flags().StringVarP(&port, "port", "p", "", "local port to start listening on")
}

// tunnelFn ...
func tunnelFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will os.Exit(1) with an error message
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
name of the container you would like to tunnel into:

ex: nanobox tunnel <container>

`, len(args))
		return
	}

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultConfig.Meta["alias"] = app
	processor.DefaultConfig.Meta["port"] = port
	processor.DefaultConfig.Meta["container"] = args[0]
	print.OutputCommandErr(processor.Run("tunnel", processor.DefaultConfig))
}
