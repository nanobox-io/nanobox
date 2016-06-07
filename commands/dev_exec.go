package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// DevExecCmd ...
	DevExecCmd = &cobra.Command{
		Use:   "exec",
		Short: "execute a command",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				fmt.Println("I need atleast one arguement to execute")
				return
			case 1:
				processor.DefaultConfig.Meta["command"] = args[0]
			default:
				processor.DefaultConfig.Meta["name"] = args[0]
				processor.DefaultConfig.Meta["command"] = strings.Join(args[1:], " ")
			}
			handleError(processor.Run("dev_console", processor.DefaultConfig))
		},
		// PostRun: halt,
	}
)
