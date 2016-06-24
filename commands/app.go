package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/app"
)

var (

	// currently a hidden command because its only used for one thing
	// AppCmd ...
	AppCmd = &cobra.Command{
		Use:   "app",
		Short: "The appelopment environment ",
		Long: `
The appelopment environment all starts below this subcommand.
		`,
		Hidden: true,
	}
)

//
func init() {
	// hidden subcommands
	AppCmd.AddCommand(app.NetfsCmd)
}
