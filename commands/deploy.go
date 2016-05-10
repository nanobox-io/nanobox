//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	DeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "get the production credentials and ensure production connection",
		Long:  ``,

		Run: func(ccmd *cobra.Command, args []string) {
			processor.Run("deploy", processor.DefaultConfig)
		},
	}

)
