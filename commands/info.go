package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/processors/sim"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// InfoCmd ...
	InfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Show information about the specified environment.",
		Long: `
Shows information about the specified environment. You must
specify which environment you would like information about.
`,
		Run: infoFn,
	}
)

// infoFn ...
func infoFn(ccmd *cobra.Command, args []string) {

	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args)

	switch location {
	case "local":
		switch name {
		case "dev":
			app, _ := models.FindAppBySlug(config.EnvID(), "dev")
			display.CommandErr(dev.Info(env, app))
		case "sim":
			app, _ := models.FindAppBySlug(config.EnvID(), "sim")
			display.CommandErr(sim.Info(env, app))
		}
	case "production":
		fmt.Println("not yet implemented")
	}
}
