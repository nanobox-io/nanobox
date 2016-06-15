package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// UpdateCmd ...
	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Updates the Nanobox CLI to the newest available *minor* version.",
		Long: `
Updates the Nanobox CLI to the newest available *minor* version.
Major version updates must be manually downloaded & installed.
		`,
		Run: updateFn,
	}
)

// updateFn ...
func updateFn(ccmd *cobra.Command, args []string) {

	// I handle this error manually because I want a different message than all
	// other commands
	if err := processor.Run("update", processor.DefaultConfig); err != nil {
		fmt.Println("Nanobox was unable to update because of the following error:\n", err.Error())
	}
}
