package sim

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/sim"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// ConsoleCmd ...
var ConsoleCmd = &cobra.Command{
	Use:    "console",
	Short:  "Opens an interactive console inside your sim platform.",
	Long:   ``,
	PreRun: steps.Run("start", "build", "compile", "sim start", "sim deploy"),
	Run:    consoleFn,
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("you need to provide a container to console into")

		return
	}

	component, _ := models.FindComponentBySlug(config.EnvID()+"_sim", args[0])

	display.CommandErr(sim.Console(component))
}
