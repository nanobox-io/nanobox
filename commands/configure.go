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
		Run:     configureFn,
		Aliases: []string{"config"},
	}

	ConfigureSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set a configuration key",
		Long: `
Set a key in the configuration		
		`,
		Run: configureSetFn,
	}

	ConfigureGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a value form the configuration",
		Long: `
Get a key from the configuration
		`,
		Run: configureGetFn,
	}

	ConfigureListCmd = &cobra.Command{
		Use:   "show",
		Short: "Show the full configuration",
		Long: `
List the full configuration.
		`,
		Run:     configureListFn,
		Aliases: []string{"list", "ls"},
	}
)

func init() {
	steps.Build("configure", configureComplete, configureFn)

	ConfigureCmd.AddCommand(ConfigureSetCmd)
	ConfigureCmd.AddCommand(ConfigureGetCmd)
	ConfigureCmd.AddCommand(ConfigureListCmd)

}

// configureFn ...
func configureFn(ccmd *cobra.Command, args []string) {

	display.CommandErr(processors.Configure())
}

func configureSetFn(ccmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("setting a key requires <key> <value>")
		return
	}
	display.CommandErr(processors.ConfigureSet(args[0], args[1]))
}

func configureGetFn(ccmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("what is the key you would like to see")
		return
	}
	config, _ := models.LoadConfig()
	jsonData, _ := json.Marshal(config)
	configMap := map[string]interface{}{}
	json.Unmarshal(jsonData, &configMap)
	fmt.Println(configMap[args[0]])
	return

}

func configureListFn(ccmd *cobra.Command, args []string) {
	config, _ := models.LoadConfig()
	prettyJson, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("%s\n", prettyJson)
	return
}

func configureComplete() bool {
	_, err := models.LoadConfig()
	return err == nil
}
