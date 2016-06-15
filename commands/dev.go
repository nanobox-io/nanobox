package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/dev"
	"github.com/nanobox-io/nanobox/processor"
	cmdutil "github.com/nanobox-io/nanobox/validate/commands"
)

var (

	// DevCmd ...
	DevCmd = &cobra.Command{
		Use:   "dev",
		Short: "Starts the Nanobox VM, provisions app, & opens an interactive terminal.",
		Long: `
Starts the Nanobox VM, provisions app, & opens an interactive
terminal. This is the primary command for managing local dev
apps and interacting with your Nanobox VM.
		`,
		PreRun: cmdutil.Validate("provider"),
		Run:    devFn,
	}
)

//
func init() {

	// public subcommands
	DevCmd.AddCommand(dev.DeployCmd)
	DevCmd.AddCommand(dev.DestroyCmd)
	DevCmd.AddCommand(dev.DNSCmd)
	DevCmd.AddCommand(dev.InfoCmd)
	DevCmd.AddCommand(dev.ExecCmd)
	DevCmd.AddCommand(dev.ConsoleCmd)
	DevCmd.AddCommand(dev.EnvCmd)
	DevCmd.AddCommand(dev.ResetCmd)

	// hidden subcommands
	DevCmd.AddCommand(dev.NetfsCmd)
}

// devFn ...
func devFn(ccmd *cobra.Command, args []string) {

	//
	if err := processor.Run("dev", processor.DefaultConfig); err != nil {

	}
}
