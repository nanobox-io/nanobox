// Package commands ...
package commands

import (
	"github.com/spf13/cobra"
)

//
var (

	//
	version bool // display the version of nanobox

	// NanoboxCmd ...
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,

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
	// NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultControl.Debug, "debug", "", false, "run nanobox in debug mode")
	// NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultControl.Force, "force", "f", false, "Forces a command to run (effects vary per command).")
	// NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultControl.Verbose, "verbose", "v", false, "Increases display output.")
	// NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultControl.Quiet, "quiet", "q", false, "Decreases display output.")

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
