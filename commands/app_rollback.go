package commands

import (
	"fmt"

	pagodaAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// AppRollbackCommand satisfies the Command interface for rolling back an app
type AppRollbackCommand struct{}

// Help prints detailed help text for the app rollback command
func (c *AppRollbackCommand) Help() {
	ui.CPrintln(`
Description:
  Rolls an application back one (1) deploy.

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

Usage:
  pagoda rollback [-a app-name]
  pagoda app:rollback [-a app-name]

  ex. pagoda rollback -a app-name

Options:
  -a, --app [app-name]
    The name of the app you want to rollback.
  `)
}

// Run attempts to roll an app back one (1) deploy
func (c *AppRollbackCommand) Run(fApp string, opts []string, api *pagodaAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	//
	if _, err := api.RollbackAppDeploy(fApp); err != nil {
		fmt.Println("There was a problem rolling back %s. See ~/.pagodabox/log.txt for details", fApp)
		ui.Error("pagoda app:rollback", err)
	}

	fmt.Printf("Rolling %s back to previous deploy... Check your dashboard for transaction details and logs\n", fApp)
}
