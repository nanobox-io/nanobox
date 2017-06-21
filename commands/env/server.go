package env

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processors/server"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// ServerCmd ...
	ServerCmd = &cobra.Command{
		Hidden: true,
		Use:    "server",
		Short:  "Server control",
		Long:   ``,
	}

	// ServerStartCmd ...
	ServerStartCmd = &cobra.Command{
		Hidden: true,
		Use:    "start",
		Short:  "Start the server",
		Long:   ``,
		Run:    serverStartFn,
	}

	// ServerStopCmd ...
	ServerStopCmd = &cobra.Command{
		Hidden: true,
		Use:    "stop",
		Short:  "Stop the server",
		Long:   ``,
		Run:    serverStopFn,
	}

	// ServerTeadownCmd ...
	ServerTeadownCmd = &cobra.Command{
		Hidden: true,
		Use:    "teardown",
		Short:  "Teardown the server",
		Long:   ``,
		Run:    serverTeardownFn,
	}
)

//
func init() {
	ServerCmd.AddCommand(ServerStartCmd)
	ServerCmd.AddCommand(ServerStopCmd)
	ServerCmd.AddCommand(ServerTeadownCmd)
}

func serverStartFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(server.Setup())
}

func serverStopFn(ccmd *cobra.Command, args []string) {

	display.CommandErr(server.Stop())
}

func serverTeardownFn(ccmd *cobra.Command, args []string) {

	display.CommandErr(server.Teardown())
}
