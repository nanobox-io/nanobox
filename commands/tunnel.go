//
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	TunnelCmd = &cobra.Command{
		Use:   "tunnel",
		Short: "do the tunnel thing",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("i need a container to run in")
				return
			}
			processor.DefaultConfig.Meta["alias"] = app
			processor.DefaultConfig.Meta["container"] = args[0]
			processor.DefaultConfig.Meta["port"] = port
			processor.Run("tunnel", processor.DefaultConfig)
		},
		// PostRun: halt,
	}

	port string
)

func init() {
	TunnelCmd.Flags().StringVarP(&app, "app", "a", "", "production app name or alias")
	TunnelCmd.Flags().StringVarP(&port, "port", "p", "", "local port to start listening on")
}
