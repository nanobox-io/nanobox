package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// BuildCmd ...
	BuildCmd = &cobra.Command{
		Use:   "build",
		Short: "do a build",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			handleError(processor.Run("build", processor.DefaultConfig))
		},
		// PostRun: halt,
	}
)
