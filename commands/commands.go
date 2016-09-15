// Package commands ...
package commands

import (
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/mixpanel"
	"github.com/nanobox-io/nanobox/util/update"
)

//
var (
	// debug mode
	debugMode bool

	// display level debug
	displayDebugMode bool

	// display level trace
	displayTraceMode bool

	//
	internalCommand bool

	// NanoboxCmd ...
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			// report the command usage to mixpanel
			mixpanel.Report(strings.Replace(ccmd.CommandPath(), "nanobox ", "", 1))

			// alert the user if an update is needed
			update.Check()

			registry.Set("internal", internalCommand)
			registry.Set("debug", debugMode)

			// setup the display output
			if displayDebugMode {
				lumber.Level(lumber.DEBUG)
				display.Summary = false
				display.Mode = "debug"
			}

			if displayTraceMode {
				lumber.Level(lumber.TRACE)
				display.Summary = false
				display.Mode = "trace"
			}

		},

		//
		Run: func(ccmd *cobra.Command, args []string) {

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

// init creates the list of available nanobox commands and sub commands
func init() {

	// persistent flags
	NanoboxCmd.PersistentFlags().BoolVarP(&internalCommand, "internal", "", false, "Skip pre-requisite checks")
	NanoboxCmd.PersistentFlags().MarkHidden("internal")
	NanoboxCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "", false, "In the event of a failure, drop into debug context")
	NanoboxCmd.PersistentFlags().BoolVarP(&displayDebugMode, "verbose", "v", false, "Increases display output and sets level to debug")
	NanoboxCmd.PersistentFlags().BoolVarP(&displayTraceMode, "trace", "t", false, "Increases display output and sets level to trace")

	// subcommands
	NanoboxCmd.AddCommand(StatusCmd)
	NanoboxCmd.AddCommand(InspectCmd)
	NanoboxCmd.AddCommand(DeployCmd)
	NanoboxCmd.AddCommand(ConsoleCmd)
	NanoboxCmd.AddCommand(LinkCmd)
	NanoboxCmd.AddCommand(LoginCmd)
	NanoboxCmd.AddCommand(LogoutCmd)
	NanoboxCmd.AddCommand(BuildCmd)
	NanoboxCmd.AddCommand(CleanCmd)
	NanoboxCmd.AddCommand(DevCmd)
	NanoboxCmd.AddCommand(SimCmd)
	NanoboxCmd.AddCommand(EnvCmd)
	NanoboxCmd.AddCommand(TunnelCmd)
	NanoboxCmd.AddCommand(DestroyCmd)
	NanoboxCmd.AddCommand(StartCmd)
	NanoboxCmd.AddCommand(StopCmd)
}
