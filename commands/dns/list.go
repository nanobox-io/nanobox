package dns

import (
	"fmt"

	"github.com/spf13/cobra"

	// "github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	app_dns "github.com/nanobox-io/nanobox/processors/app/dns"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// ListCmd ...
var ListCmd = &cobra.Command{
	Use:   "ls [local|dry-run]",
	Short: "list dns entries",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: listFn,
}

// listFn ...
func listFn(ccmd *cobra.Command, args []string) {

	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 0)

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_dns.List(app))
	case "production":
		fmt.Println("not yet implemented")
	}
}
