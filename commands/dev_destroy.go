package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// DevDestroyCmd ...
	DevDestroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroys the docker machines associated with the current app.",
		Long:  `
Destroys the docker machines associated with the current app.
If no other apps are running, it will destroy the Nanobox VM.
		`,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			handleError(processor.Run("dev_destroy", processor.DefaultConfig))
		},
		// PostRun: halt,
	}
)
