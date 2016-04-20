// package commands ...
package commands

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	_ "github.com/nanobox-io/nanobox/processor/code"
	_ "github.com/nanobox-io/nanobox/processor/nanopack"
	_ "github.com/nanobox-io/nanobox/processor/provider"
	_ "github.com/nanobox-io/nanobox/processor/service"
	"github.com/nanobox-io/nanobox/validate"
)

const VERSION = "1.0.0"

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
			if processor.DefaultConfig.Verbose {
				lumber.Level(lumber.DEBUG)
			}
		},

		//
		Run: func(ccmd *cobra.Command, args []string) {

			// hijack the verbose flag (-v), and use it to display the version of the
			// CLI
			if version || processor.DefaultConfig.Verbose {
				fmt.Printf("nanobox v%s\n", version)
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
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultConfig.DevMode, "dev", "", false, "")
	NanoboxCmd.PersistentFlags().MarkHidden("dev")

	// persistent flags
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultConfig.Background, "background", "", false, "Stops nanobox from auto-suspending.")
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultConfig.Force, "force", "f", false, "Forces a command to run (effects vary per command).")
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultConfig.Verbose, "verbose", "v", false, "Increase command output from 'info' to 'debug'.")

	// local flags
	NanoboxCmd.Flags().BoolVarP(&version, "version", "", false, "Display the current version of this CLI")

	// nanobox commands
	// comment to get things to work .. will be implemented later
	// NanoboxCmd.AddCommand(updateCmd)

	// subcommands
	NanoboxCmd.AddCommand(DevCmd)
	NanoboxCmd.AddCommand(BuildCmd)
}

func validCheck(checks ...string) func(ccmd *cobra.Command, args []string) {
	return func(ccmd *cobra.Command, args []string) {
		if err := validate.Check(checks...); err != nil {
			fmt.Printf("Validation Failed:\n%s\n", err.Error())
			os.Exit(1)
		}
	}
}
