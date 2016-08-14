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

	// InfoCmd ...
	InfoCmd = &cobra.Command{
		Use:    "info",
		Short:  "Displays information about the running sim app and its components.",
		Long:   ``,
		PreRun: validate.Requires("provider"),
		Run:    infoFn,
	}
)

// infoFn will run the DNS processor for adding DNS entires to the "hosts" file
func infoFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	simInfo := sim.Info{App: app}
	print.OutputCommandErr(simInfo.Run())
}
