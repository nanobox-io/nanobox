package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
)

var (
	username string
	password string

	// LoginCmd ...
	LoginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticates your nanobox client with your nanobox.io account.",
		Long:  ``,
		Run:   loginFn,
	}
)

//
func init() {
	LoginCmd.Flags().StringVarP(&username, "username", "u", "", "username")
	LoginCmd.Flags().StringVarP(&password, "password", "p", "", "password")
}

// loginFn ...
func loginFn(ccmd *cobra.Command, args []string) {
	processor.DefaultConfig.Meta["username"] = username
	processor.DefaultConfig.Meta["password"] = password
	print.OutputCommandErr(processor.Run("login", processor.DefaultConfig))
}
