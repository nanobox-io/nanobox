package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
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
	// if they want to configure a key/value
	// show the config
	if len(args) == 1 {
		config, _ := models.LoadConfig()
		prettyJson, _ := json.MarshalIndent(config, "", "  ")
		fmt.Printf("%s\n", prettyJson)
		return

	}
	if len(args) == 2 {
		display.CommandErr(processors.ConfigureSet(args[0], args[1]))
		return
	}
	display.CommandErr(processors.Configure())
}

func configureComplete() bool {
	_, err := models.LoadConfig()
	return err == nil
}
