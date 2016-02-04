// package commands ...
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/box"
	"github.com/nanobox-io/nanobox/commands/dev"
	"github.com/nanobox-io/nanobox/commands/engine"
	"github.com/nanobox-io/nanobox/commands/service"
	"github.com/nanobox-io/nanobox/config"
)

//
var (

	//
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,

		// if the verbose flag is used, up the log level (this will cascade into
		// all subcommands since this is the root command)
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			if config.Verbose {
				config.LogLevel = "debug"
			}
		},

		//
		Run: func(ccmd *cobra.Command, args []string) {

			// hijack the verbose flag (-v), and use it to display the version of the
			// CLI
			if version || config.Verbose {
				fmt.Printf("nanobox v%s\n", config.VERSION)
				return
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}

	//
	version bool // display the version of nanobox
)

// init creates the list of available nanobox commands and sub commands
func init() {

	// internal flags
	NanoboxCmd.PersistentFlags().BoolVarP(&config.Devmode, "dev", "", false, "")
	NanoboxCmd.PersistentFlags().MarkHidden("dev")

	// persistent flags
	NanoboxCmd.PersistentFlags().BoolVarP(&config.Background, "background", "", false, "Stops nanobox from auto-suspending.")
	NanoboxCmd.PersistentFlags().BoolVarP(&config.Force, "force", "f", false, "Forces a command to run (effects vary per command).")
	NanoboxCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Increase command output from 'info' to 'debug'.")

	// local flags
	NanoboxCmd.Flags().BoolVarP(&version, "version", "", false, "Display the current version of this CLI")

	// nanobox commands
	NanoboxCmd.AddCommand(publishCmd)
	NanoboxCmd.AddCommand(updateCmd)

	// subcommands
	NanoboxCmd.AddCommand(box.BoxCmd)
	NanoboxCmd.AddCommand(dev.DevCmd)
	NanoboxCmd.AddCommand(engine.EngineCmd)
	NanoboxCmd.AddCommand(service.ServiceCmd)
}
