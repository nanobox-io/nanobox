package commands

import (
	"fmt"
	"os"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// AppInfoCommand satisfies the Command interface for obtaining app info
type AppInfoCommand struct{}

// Help prints detailed help text for the app info command
func (c *AppInfoCommand) Help() {
	ui.CPrintln(`
Description:
  Display information about an application.

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

Usage:
  pagoda info [-a app-name]
  pagoda app:info [-a app-name]

  ex. pagoda info -a app-name

Options:
  -a, --app [app-name]
    The name of the app you want information for.
  `)
}

// Run prints out select information for the designated app
func (c *AppInfoCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	//
	app, err := api.GetApp(fApp)
	if err != nil {
		fmt.Printf("Unable to find an app with the name '%v' \n", fApp)
		os.Exit(1)
	}

	appGitURL := "git@git.pagodabox.io:apps/" + app.Name + ".git"
	appFlation := helpers.DetermineAppFlation(app.Flation)
	appColor := helpers.DetermineAppStatus(app.State, app.Flation)
	appType := helpers.DetermineAppType(app.Free)

	var appName string

	switch app.State {
	case "initialized", "created", "active":
		appName = appColor + appType + app.Name + "[reset]"
	case "uninitialized", "inactive", "defunct":
		appName = appColor + "x " + appType + app.Name + "[reset]"
	case "hybernated":
		appName = appColor + "! " + appType + app.Name + "[reset]"
	}

	// general info
	ui.CPrintln(`
General Information:
-------------------------
App Name  : ` + appName + `
Clone URL : ` + appGitURL + `
Status    : ` + appFlation + `
Timezone  : ` + app.Timezone + `
  `)

	// deploy info
	fmt.Println(`
Currently Deployed:
-------------------------`)

	if app.ActiveDeployID != "" {
		deploy, err := api.GetAppDeploy(app.ID, app.ActiveDeployID)
		if err != nil {
			fmt.Printf("Oops! We could not find any deploys for '%v'.\n", app.Name)
			os.Exit(1)
		}

		fmt.Println("Commit  : " + deploy.Commit)
		fmt.Println("Message : " + deploy.Message)
	} else {
		fmt.Println("No code has been deployed to this app.")
	}

	fmt.Println("")

}
