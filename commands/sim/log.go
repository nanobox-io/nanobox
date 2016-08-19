package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/sim"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// LogCmd ...
	LogCmd = &cobra.Command{
		Use:    "log",
		Short:  "Displays logs from the running sim app and its components.",
		Long:   ``,
		PreRun: validate.Requires("provider", "provider_up", "built"),
		Run:    logFn,
	}
)

// logFn will run the DNS processor for adding DNS entires to the "hosts" file
func logFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	simLog := sim.Log{
		App: app,
	}
	display.CommandErr(simLog.Run())
}
