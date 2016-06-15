package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	cmdutil "github.com/nanobox-io/nanobox/validate/commands"
)

var (

	// ConsoleCmd ...
	ConsoleCmd = &cobra.Command{
		Use:    "console",
		Short:  "Opens an interactive console inside a production component.",
		Long:   ``,
		PreRun: cmdutil.Validate("provider"),
		Run:    consoleFn,
	}
)

//
func init() {
	ConsoleCmd.Flags().StringVarP(&app, "app", "a", "", "app-name or alias")
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("I need a component to run in")
		return
	}
	processor.DefaultConfig.Meta["alias"] = app
	processor.DefaultConfig.Meta["container"] = args[0]

	//
	if err := processor.Run("console", processor.DefaultConfig); err != nil {

	}
}
