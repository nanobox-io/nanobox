package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// LoginCmd ...
	LoginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticates your nanobox client with your nanobox.io account.",
		Long:  ``,

		Run: func(ccmd *cobra.Command, args []string) {
			processor.DefaultConfig.Meta["username"] = username
			processor.DefaultConfig.Meta["password"] = password
			handleError(processor.Run("login", processor.DefaultConfig))
		},
	}

	username string
	password string
)

//
func init() {
	LoginCmd.Flags().StringVarP(&username, "username", "u", "", "username")
	LoginCmd.Flags().StringVarP(&password, "password", "p", "", "password")
}
