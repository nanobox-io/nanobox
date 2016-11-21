package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/platform"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// LogCmd ...
var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "View and streams application logs.",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: logFn,
}

// logFn ...
func logFn(ccmd *cobra.Command, args []string) {

	// parse the evars excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 1)

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(platform.MistListen(app))
	case "production":
		fmt.Println("not yet implemented")
	}
}
