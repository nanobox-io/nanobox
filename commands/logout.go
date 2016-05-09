//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	LogoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "remove api token and forget about the user",
		Long:  ``,

		Run: func(ccmd *cobra.Command, args []string) {
			processor.Run("logout", processor.DefaultConfig)
		},
	}
)
