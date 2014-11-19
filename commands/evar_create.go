package commands

import (
	"fmt"
	"os"
	"regexp"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// EVarCreateCommand satisfies the Command interface for creating an environment
// variable
type EVarCreateCommand struct{}

// Help prints detailed help text for the evar create command
func (c *EVarCreateCommand) Help() {
	ui.CPrintln(`
Description:
  Creates a new environment variable for [app-name].

  If [app-name] is not specified, a name will be generated for you.

Usage:
  pagoda evar:create [-a app-name] KEY=VALUE

  ex. pagoda evar:create -a app-name ENV=PRODUCTION

Options:
  -a, --app [app-name]
    The name of the app
  `)
}

// Run attempts to create a new app on Nanobox. It can take an app-name flag
// for naming the app, and a tinker flag for designating the type of app to create.
// If successful, it attempts to add a new remote, then prints instructions on
// pushing code to pagodabox
func (c *EVarCreateCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	// if fApp == "" {
	// 	fApp = helpers.FindPagodaApp()
	// }
	fApp = "TESTING"

	var fEVar string

	// If there's no environment variable, message and exit
	if len(opts) <= 0 {
		fmt.Printf(`
Oops! You forgot include an environment variable:
  ex. pagoda evar:create KEY=VALUE
    `)
		os.Exit(1)
	} else {
		fEVar = opts[0]
	}

	reFindEVar := regexp.MustCompile(`^(.+)\=(.+)$`)

	subMatch := reFindEVar.FindStringSubmatch(fEVar)
	if subMatch == nil {
		fmt.Printf("Your environment variable '%v' does not match the required format 'KEY=VALUE'", fEVar)
		os.Exit(1)
	}

	//
	eVarCreateOptions := &nanoAPI.EVarCreateOptions{Title: subMatch[1], Value: subMatch[2]}

	// create evar
	eVar, err := api.CreateEVar(eVarCreateOptions)
	if err != nil {
		_, err, msg := helpers.HandleAPIError(err)
		fmt.Printf("Oops! We could not create your evar: %v - %v", err, msg)
		os.Exit(1)
	}

	ui.CPrintln(`
New environment variable added to [green]` + fApp + `[reset]:
  ` + eVar.Title + ` = ` + eVar.Value)
}
