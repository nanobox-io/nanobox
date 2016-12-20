
package env

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processors/provider/bridge"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// BridgeCmd ...
	BridgeCmd = &cobra.Command{
		Hidden: true,
		Use:    "bridge",
		Short:  "Bridge control",
		Long:   ``,
	}

	// BridgeStartCmd ...
	BridgeStartCmd = &cobra.Command{
		Hidden: true,
		Use:    "start",
		Short:  "Start the bridge",
		Long:   ``,
		Run:    bridgeStartFn,
	}

	// BridgeStopCmd ...
	BridgeStopCmd = &cobra.Command{
		Hidden: true,
		Use:    "stop",
		Short:  "Stop the bridge",
		Long:   ``,
		Run:    bridgeStopFn,
	}

	// BridgeTeadownCmd ...
	BridgeTeadownCmd = &cobra.Command{
		Hidden: true,
		Use:    "teardown",
		Short:  "Teardown the bridge",
		Long:   ``,
		Run:    bridgeTeardownFn,
	}
)

//
func init() {
	BridgeCmd.AddCommand(BridgeStartCmd)
	BridgeCmd.AddCommand(BridgeStopCmd)
	BridgeCmd.AddCommand(BridgeTeadownCmd)
}

func bridgeStartFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(bridge.Start())
}

func bridgeStopFn(ccmd *cobra.Command, args []string) {

	display.CommandErr(bridge.Stop())
}

func bridgeTeardownFn(ccmd *cobra.Command, args []string) {

	display.CommandErr(bridge.Teardown())
}
