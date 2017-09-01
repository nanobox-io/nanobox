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

var (
	logFollow bool
	logNumber int
	logRaw    bool   // display log timestamps instead of added ones
	logStart  string // todo: forthcoming
	logEnd    string // todo: forthcoming
	logLimit  string // todo: forthcoming

	// LogCmd provides the logging functionality.
	LogCmd = &cobra.Command{
		Use:   "log [dry-run|remote-alias]",
		Short: "Streams application logs.",
		Long:  "'remote-alias' is the alias for your app, given on `nanobox remote add app-name alias`",
		Run:   logFn,
	}
)

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
		logOpts := models.LogOpts{
			Number: logNumber,
			Follow: logFollow,
			Raw:    logRaw,
		}

		// since we default to live logging, if `-n` is set, we'll print that many
		// historic logs and return unless `-f` is also set.
		if logNumber > 0 {
			display.CommandErr(log.Print(envModel, name, logOpts))

			// if `-f` is also specified, continue, else return here.
			if !logFollow {
				return
			}
		}

		// set the meta arguments to be used in the processor and run the processor
		display.CommandErr(log.Tail(envModel, name, logOpts))
	}
}
