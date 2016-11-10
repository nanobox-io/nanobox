package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/evar"
)

var (

	// EvarCmd ...
	EvarCmd = &cobra.Command{
		Use:   "evar",
		Short: "Manage environment variables.",
		Long: `
Manages environment variables in your different environments.
		`,
	}
)

//
func init() {
	EvarCmd.AddCommand(evar.AddCmd)
	EvarCmd.AddCommand(evar.RemoveCmd)
	EvarCmd.AddCommand(evar.ListCmd)
}
