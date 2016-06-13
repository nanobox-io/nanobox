package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// DevDeployCmd ...
	DevDeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys your build package into your Nanobox VM and starts all services.",
		Long:  `
Deploys your build package into your Nanobox VM and
starts all services. This is used to simulate a full
deploy locally, before deploying into production.
		`,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			handleError(processor.Run("dev_deploy", processor.DefaultConfig))
		},
		// PostRun: halt,
	}
)
