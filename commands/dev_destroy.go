//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	DevDestroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "destroy the docker machines associated witht his app",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			processor.Run("dev_destroy", processor.DefaultConfig)
		},
		// PostRun: halt,
	}
)
