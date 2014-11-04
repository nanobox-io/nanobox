package commands

import (
	"fmt"
	"os"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// EVarDestroyCommand satisfies the Command interface for destroying an app's
// environemtn variable
type EVarDestroyCommand struct{}

// Help prints detailed help text for the evar destroy command
func (c *EVarDestroyCommand) Help() {
	ui.CPrintln(`
Description:
  Destroys an environment variable

  If [app-name] is not specified, will prompt for one. Also, before an app is
  destroyed, will prompt for [app-name] to confirm action.

Usage:
  pagoda evar:destroy [-a app-name] evar-key

  ex. pagoda evar:destroy -a app-name env

Options:
  -a, --app [app-name]
    The name of the app
  `)
}

// Run attempts to destroy an app on Nanobox. It can take a force flag that will
// skip the confirmation process, other wise will ask for confirmation by retyping
// the name of the app to be destroyed
func (c *EVarDestroyCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	var fEVar string

	// get environment variables
	eVars, err := api.GetAppEVars(fApp)
	if err != nil {
		fmt.Printf("There was a problem getting '%s's' environment variables. See ~/.pagodabox/log.txt for details", fApp)
		ui.Error("pagoda evar:destroy", err)
	}

	// If there's no key, prompt for one
	if len(opts) <= 0 {
		fmt.Println(`
We found the following environment variables tied to this app:

KEY = VALUE
--------------------------------------------------`)
		for _, eVar := range eVars {
			if !eVar.Internal {
				fmt.Printf("- %s = %s", eVar.Title, eVar.Value)
				fmt.Println("")
			}
		}
		fmt.Println("")
		fEVar = ui.Prompt("Which environment variable (by KEY) would you like to destroy: ")

		// We should expect opts[0] to be the key.
	} else {
		fEVar = opts[0]
	}

	//
	var eVarID string

	for _, eVar := range eVars {
		if fEVar == eVar.Title {
			eVarID = eVar.ID
		}
	}

	// destroy evar
	if err := api.DeleteAppEVar(fApp, eVarID); err != nil {
		_, err, msg := helpers.HandleAPIError(err)
		fmt.Printf("Oops! We could not destroy your environment variable: %s - %s", err, msg)
		os.Exit(1)
	}

	fmt.Printf("Environment variable '%s' destroyed for '%s'", fEVar, fApp)

}
