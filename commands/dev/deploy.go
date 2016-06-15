package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	cmdutil "github.com/nanobox-io/nanobox/validate/commands"
)

// DeployCmd ...
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys your build package into your Nanobox VM and starts all services.",
	Long: `
Deploys your build package into your Nanobox VM and
starts all services. This is used to simulate a full
deploy locally, before deploying into production.
		`,
	PreRun: cmdutil.Validate("provider"),
	Run:    deployFn,
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {

	//
	if err := processor.Run("dev_deploy", processor.DefaultConfig); err != nil {

	}
}
