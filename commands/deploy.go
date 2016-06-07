package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// DeployCmd ...
	DeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "get the production credentials and ensure production connection",
		Long:  ``,

		Run: func(ccmd *cobra.Command, args []string) {
			processor.DefaultConfig.Meta["alias"] = "default"
			if len(args) == 1 {
				processor.DefaultConfig.Meta["alias"] = args[0]
			}
			handleError(processor.Run("deploy", processor.DefaultConfig))
		},
	}
)
