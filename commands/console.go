//
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	ConsoleCmd = &cobra.Command{
		Use:   "console",
		Short: "do the console thing",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("i need a container to run in")
				return
			}
			processor.DefaultConfig.Meta["alias"] = app
			processor.DefaultConfig.Meta["container"] = args[0]
			processor.Run("console", processor.DefaultConfig)
		},
		// PostRun: halt,
	}
)

func init() {
	ConsoleCmd.Flags().StringVarP(&app, "app", "a", "", "production app name or alias")
}
