package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/log"
	"github.com/nanobox-io/nanobox/processors/platform"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// LogCmd ...
var LogCmd = &cobra.Command{
	Use:    "log [dry-run|remote-alias]",
	Short:  "Streams application logs.",
	Long:   "'remote-alias' is the alias for your app, given on `nanobox remote add app-name alias`",
	Run:    logFn,
}

// logFn ...
func logFn(ccmd *cobra.Command, args []string) {

	// parse the evars excluding the context
	envModel, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(envModel, args, 1)

	switch location {
	case "local":
		if name == "dev" {
			fmt.Printf(`
--------------------------------------------------------
Watching 'local' not yet implemented. You can watch your
logs inside a terminal running 'nanobox run'.
--------------------------------------------------------

`)
			return
		}
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(platform.MistListen(app))
	case "production":
		steps.Run("login")(ccmd, args)

		// set the meta arguments to be used in the processor and run the processor
		display.CommandErr(log.Tail(envModel, name))
	}
}
