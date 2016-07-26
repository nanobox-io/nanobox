package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

// DeployCmd ...
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a build package into your sim platform and starts all services.",
	Long: `
Deploys a build package into your sim platform and
starts all services. This is used to simulate a full
deploy locally, before deploying into production.
		`,
	PreRun: validate.Requires("provider", "provider_up", "built"),
	Run:    deployFn,
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("sim_deploy", processor.DefaultControl))
}
