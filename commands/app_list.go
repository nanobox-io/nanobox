package commands

import (
	"fmt"

	pagodaAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// AppListCommand satisfies the Command interface for listing a user's apps
type AppListCommand struct{}

// Help prints detailed help text for the app list command
func (c *AppListCommand) Help() {
	ui.CPrintln(`
Description:
  Lists all of your applications.

  state:
    [blue]created[reset]    - App exists, but has not code deployed
    [green]active[reset]     - App exists, and has been deployed
    [yellow]asleep[reset]     - A sleeping 'tinker' app
    [red]inactive[reset]   - App is queued to be deleted
    [red]hybernated[reset] - App is disabled due to deliquency

  name:
    The name's of your applications.

  permission:
    Your access level in relation to the app.

    owner     : Full permissions
    manager   : CANNOT delete the app or modify billing info
    developer : CAN push, pull, and deploy code only.

  type:
    'Tinker' apps are preceded by a '*'.

Usage:
  pagoda list
  pagoda app:list

  ex. pagoda list
  `)
}

// Run displays select information about all of a user's apps
func (c *AppListCommand) Run(fApp string, opts []string, api *pagodaAPI.Client) {

	// get apps
	apps, err := api.GetApps()
	if err != nil {
		fmt.Println("There was a problem getting your apps. See ~/.pagodabox/log.txt for details\n")
		ui.Error("pagoda app:list", err)
	}

	// get roles
	roles, err := api.GetUserRoles()
	if err != nil {
		fmt.Println("There was a problem getting your apps role's. See ~/.pagodabox/log.txt for details\n")
		ui.Error("pagoda app:list", err)
	}

	rolesMap := make(map[string]string)
	for _, role := range roles {
		rolesMap[role.AppID] = role.Permission
	}

	//
	fmt.Println(`
state name (permission)
------------------------------`)

	var appColor, appType string

	for _, app := range apps {

		appColor = helpers.DetermineAppStatus(app.State, app.Flation)
		appType = helpers.DetermineAppType(app.Free)

		switch app.State {

		//
		case "initialized", "created", "active":
			ui.CPrint(appColor + "\u25CF[reset] " + appType + app.Name)

		//
		case "uninitialized", "inactive", "defunct":
			ui.CPrint(appColor + "x[reset] " + appType + app.Name)

		//
		case "hybernated":
			ui.CPrint(appColor + "![reset] " + appType + app.Name)
		}

		if val, ok := rolesMap[app.ID]; ok {
			fmt.Print(" (" + val + ")")
		}

		fmt.Println("")
	}

	fmt.Println("")
}
