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
	if len(args) != 1 {
		fmt.Println("i need a container to run in")
		return
	}

	//
	processor.DefaultConfig.Meta["alias"] = app
	processor.DefaultConfig.Meta["container"] = args[0]
	processor.DefaultConfig.Meta["port"] = port
	print.OutputCommandErr(processor.Run("tunnel", processor.DefaultConfig))
}
