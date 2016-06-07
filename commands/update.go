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
		Short: "Updates this CLI to the newest available version (does not update major versions)",
		Long:  ``,

		Run: func(ccmd *cobra.Command, args []string) {

			// I want to handle this error manually because I want a special message;
			//
			// NOTE: handleError could be updated at some point to accept a msg string,
			// I wont do that now because I don't know the entire scope yet.
			if err := processor.Run("update", processor.DefaultConfig); err != nil {
				fmt.Println("Nanobox was unable to update because of the following error:\n", err.Error())
			}
		},
	}
)
