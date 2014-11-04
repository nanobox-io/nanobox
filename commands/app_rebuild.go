package commands

import (
	"fmt"

	pagodaAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// AppRebuildCommand satisfies the Command interface for rebuilding an app
type AppRebuildCommand struct{}

// Help prints detailed help text for the app rebuild command
func (c *AppRebuildCommand) Help() {
	ui.CPrintln(`
Description:
  Rebuilds an application's libraries and redeploys all code services.

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

Usage:
  pagoda rebuild [-a app-name]
  pagoda app:rebuild [-a app-name]

  ex. pagoda rebuild -a app-name

Options:
  -a, --app [app-name]
    The name of the app you want to rebuild.
  `)
}

// Run attempts to rebuild an app
func (c *AppRebuildCommand) Run(fApp string, opts []string, api *pagodaAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	//
	if _, err := api.RebuildAppLibs(fApp); err != nil {
		fmt.Println("There was a problem rebuilding %s. See ~/.pagodabox/log.txt for details", fApp)
		ui.Error("pagoda app:rebuild", err)
	}

	fmt.Printf("Rebuilding %s... Check your dashboard for transaction details and logs\n", fApp)
}
