package dns

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/server"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	app_dns "github.com/nanobox-io/nanobox/processors/app/dns"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)


func init() {
	server.AddCmd("dns add", addFn)	
}

// AddCmd ...
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds dns entries",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: addFn,
}

// addFn ...
func addFn(ccmd *cobra.Command, args []string) {

	// parse the dnss excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 2)

	if len(args) != 1 {
		fmt.Println("i need a dns")
	}

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		app.Generate(env, name)
		display.CommandErr(app_dns.Add(env, app, args[0]))
	case "production":
		fmt.Println("not yet implemented")
	}
}
