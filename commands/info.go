package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// InfoCmd ...
	InfoCmd = &cobra.Command{
		Use:   "info [local | dry-run]",
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
	args, location, name := helpers.Endpoint(env, args, 0)

	switch location {
	case "local":
		appModel, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app.Info(env, appModel))
	case "production":
		fmt.Printf(`
----------------------------------------------------------
Showing production app information is not yet implemneted.
----------------------------------------------------------

`)
	}
}
