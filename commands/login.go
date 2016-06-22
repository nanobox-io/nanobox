package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
)

var (

	// LoginCmd ...
	LoginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticates your nanobox client with your nanobox.io account.",
		Long:  ``,
		Run:   loginFn,
	}

	// loginCmdFlags ...
	loginCmdFlags = struct {
		username string
		password string
	}{}
)

//
func init() {
	LoginCmd.Flags().StringVarP(&loginCmdFlags.username, "username", "u", "", "username")
	LoginCmd.Flags().StringVarP(&loginCmdFlags.password, "password", "p", "", "password")
}

// loginFn ...
func loginFn(ccmd *cobra.Command, args []string) {

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultControl.Meta["username"] = loginCmdFlags.username
	processor.DefaultControl.Meta["password"] = loginCmdFlags.password
	print.OutputCommandErr(processor.Run("login", processor.DefaultControl))
}
