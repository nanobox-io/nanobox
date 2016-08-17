// Package commands ...
package commands

import (
	"github.com/spf13/cobra"
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/display"
)

//
var (

	// display level debug
	displayDebugMode bool
	
	// display level trace
	displayTraceMode bool

	// NanoboxCmd ...
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {

			// setup the display output
			if displayDebugMode {
				lumber.Level(lumber.DEBUG)
				display.Summary = false
				display.Mode    = "debug"
			}
			
			if displayTraceMode {
				lumber.Level(lumber.TRACE)
				display.Summary = false
				display.Mode    = "trace"
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

	// commented because this part is changing
	// // persistent flags
	NanoboxCmd.PersistentFlags().BoolVarP(&displayDebugMode, "v", "", false, "Increases display output and sets level to debug")
	NanoboxCmd.PersistentFlags().BoolVarP(&displayTraceMode, "vv", "", false, "Increases display output and sets level to trace")

	// // local flags
	// NanoboxCmd.Flags().BoolVarP(&version, "version", "", false, "Displays the current version of this CLI.")

	// nanobox commands
	NanoboxCmd.AddCommand(UpdateCmd)

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
	NanoboxCmd.AddCommand(DestroyCmd)
	NanoboxCmd.AddCommand(StartCmd)
	NanoboxCmd.AddCommand(StopCmd)

}
