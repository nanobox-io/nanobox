package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// ConfigureCmd ...
	ConfigureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Configure Nanobox.",
		Long: `
Walks through a series of question prompts that modify your local
Nanobox configuration (~/.nanobox/config.yml).
		`,
		Run: configureFn,
	}
)

func init() {
	steps.Build("configure", configureComplete, configureFn)
}

// configureFn ...
func configureFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Configure())
}

func configureComplete() bool {
	_, err := models.LoadConfig()
	return err == nil
}
