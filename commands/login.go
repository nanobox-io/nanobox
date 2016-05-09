//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	LoginCmd = &cobra.Command{
		Use:   "login",
		Short: "get the production credentials and ensure production connection",
		Long:  ``,

		Run: func(ccmd *cobra.Command, args []string) {
			processor.DefaultConfig.Meta["username"] = username
			processor.DefaultConfig.Meta["password"] = password
			processor.Run("login", processor.DefaultConfig)
		},
	}

	username string
	password string
)

func init() {
	LoginCmd.Flags().StringVarP(&username, "username", "u", "", "username")
	LoginCmd.Flags().StringVarP(&password, "password", "p", "", "password")
}
