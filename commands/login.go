package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/commands/steps"
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

	steps.Build("login", loginCheck, loginFn)
}

// loginFn ...
func loginFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Login(loginCmdFlags.username, loginCmdFlags.password))
}

func loginCheck() bool {
	auth, _ := models.LoadAuth()
	return auth.Key != ""
}
