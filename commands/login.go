package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// LoginCmd ...
	LoginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticate your nanobox client with your nanobox.io account.",
		Long:  `
Authenticate with your nanobox account by passing the username and password in or using the following environment variables:
NANOBOX_USERNAME
NANOBOX_PASSWORD
		`,
		Run:   loginFn,
	}

	// loginCmdFlags ...
	loginCmdFlags = struct {
		username string
		password string
		endpoint string
	}{}
)

//
func init() {
	LoginCmd.Flags().StringVarP(&loginCmdFlags.username, "username", "u", "", "username")
	LoginCmd.Flags().StringVarP(&loginCmdFlags.password, "password", "p", "", "password")
	LoginCmd.Flags().StringVarP(&loginCmdFlags.endpoint, "endpoint", "e", "", "endpoint")

	steps.Build("login", loginCheck, loginFn)
}

// loginFn ...
func loginFn(ccmd *cobra.Command, args []string) {
	err := processors.Login(loginCmdFlags.username, loginCmdFlags.password, loginCmdFlags.endpoint)

	display.CommandErr(err)
}

func loginCheck() bool {
	auth, _ := models.LoadAuth()
	return auth.Key != ""
}
