package sim

import (
	// "fmt"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

func init() {
	steps.Build("sim stop", true, stopCheck, stopFn)
}

// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	appModel, _ := models.FindAppBySlug(config.EnvID(), "sim")
	display.CommandErr(app.Stop(appModel))
}

func stopCheck() bool {
	// currently we always stop if we are asking weather to stop
	return false
}
