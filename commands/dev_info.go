package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// DevInfoCmd ...
	DevInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Displays information about the running app and its components.",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			handleError(processor.Run("dev_info", processor.DefaultConfig))
		},
		// PostRun: halt,
	}
)
