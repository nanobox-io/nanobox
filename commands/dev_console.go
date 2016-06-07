package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// DevConsoleCmd ...
	DevConsoleCmd = &cobra.Command{
		Use:   "console",
		Short: "do the console thing",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("i need a container to run in")
				return
			}
			processor.DefaultConfig.Meta["name"] = args[0]
			handleError(processor.Run("dev_console", processor.DefaultConfig))
		},
		// PostRun: halt,
	}
)
