package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// DevCmd ...
	DevCmd = &cobra.Command{
		Use:   "dev",
		Short: "Starts the Nanobox VM, provisions app, & opens an interactive terminal.",
		Long:  `
Starts the Nanobox VM, provisions app, & opens an interactive
terminal. This is the primary command for managing local dev
apps and interacting with your Nanobox VM.
		`,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			handleError(processor.Run("dev", processor.DefaultConfig))
		},
		// PostRun: halt,
	}
)

//
func init() {
	// public subcommands
	DevCmd.AddCommand(DevDeployCmd)
	DevCmd.AddCommand(DevDestroyCmd)
	DevCmd.AddCommand(DevInfoCmd)
	DevCmd.AddCommand(DevExecCmd)
	DevCmd.AddCommand(DevConsoleCmd)
	DevCmd.AddCommand(DevEnvCmd)
	DevCmd.AddCommand(DevResetCmd)

	// hidden subcommands
	DevCmd.AddCommand(DevNetfsCmd)
}
