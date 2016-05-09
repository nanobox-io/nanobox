//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	InfoCmd = &cobra.Command{
		Use:   "info",
		Short: "get this things information",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			processor.Run("dev_info", processor.DefaultConfig)
		},
		// PostRun: halt,
	}
)
