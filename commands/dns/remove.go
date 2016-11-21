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

// RemoveCmd ...
var RemoveCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove environment variable(s)",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: removeFn,
}

// removeFn ...
func removeFn(ccmd *cobra.Command, args []string) {
	// parse the dnss excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 0)

	if len(args) != 1 {
		fmt.Println("i need a dns")
	}

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_dns.Remove(app, args[0]))
	case "production":
		fmt.Println("not yet implemented")
	}
}
