//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	DevDeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "start a nanobox application as if it is in production",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			processor.Run("dev_deploy", processor.DefaultConfig)
		},
		// PostRun: halt,
	}
)
