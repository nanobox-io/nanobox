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

// RemoveAllCmd ...
var RemoveAllCmd = &cobra.Command{
	Use:    "rm-all",
	Short:  "remove all the dns vars",
	Long:   ``,
	Run:    removeAllFn,
	Hidden: true,
}

// removeAllFn ...
func removeAllFn(ccmd *cobra.Command, args []string) {
	// parse the dnss excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	_, location, name := helpers.Endpoint(env, args)

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_dns.RemoveAll(app))
	case "production":
		fmt.Println("not yet implemented")
	}
}
