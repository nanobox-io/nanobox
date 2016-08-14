package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/sim"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// StartCmd ...
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts your sim platform.",
		Long: `
Starts the sim platform from its previous state. If starting for
the first time, you should also generate a build (nanobox build)
and deploy it into your sim platform (nanobox sim deploy).
		`,
		PreRun: validate.Requires("provider"),
		Run:    startFn,
	}
)

// startFn ...
func startFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")

	simStart := sim.Start{App: app}
	print.OutputCommandErr(simStart.Run())
}
