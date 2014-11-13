package commands

import (
	"fmt"
	"os"
	"strconv"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// ServiceInfoCommand satisfies the Command interface for obtaining service info
type ServiceInfoCommand struct{}

// Help prints detailed help text for the service info command
func (c *ServiceInfoCommand) Help() {
	ui.CPrintln(`
Description:
  Displays info for a service in an app

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

  The service-name/UID is [yellow]required[reset]. (Ex. web1)

Usage:
  pagoda service:info [-a app-name] service-name/UID

  ex. pagoda service:info -a app-name web1

Options:
  -a, --app [app-name]
    The name of the app
  `)
}

// Run attempts to print select information about an app's service. If no service
// is provided the command will prompt for one
func (c *ServiceInfoCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	var fService string

	// If there's no service, prompt for one
	if len(opts) <= 0 {
		fService = ui.Prompt("Which service would you like see information for: ")

		// We should expect opts[0] to be the service.
	} else {
		fService = opts[0]
	}

	//
	service, err := helpers.GetServiceBySlug(fApp, fService, api)
	if err != nil {
		fmt.Printf("Oops! We could not find a '%v' on '%v'.\n", fService, fApp)
		os.Exit(1)
	}

	//
	fmt.Println(`
General Information:
------------------------------
Name      : ` + service.Name + `
UID       : ` + service.UID + `
State     : ` + service.State + `
Topology  : ` + service.Topology)

	// only show connection details for 'authable' services
	if service.Authable {
		fmt.Println(`

Connection:
------------------------------
Host      : ` + service.TunnelIP + `
Username  : ` + service.Usernames["default"] + `
Password  : ` + service.Passwords["default"] + `
Database  : ` + service.TunnelUser)
	}

	fmt.Println(`

SSH Credentials:
------------------------------
Host      : ` + service.TunnelIP + `
Port      : ` + strconv.Itoa(service.TunnelPort) + `
User      : ` + service.TunnelUser + `
Password  : (Public SSH Key)
  `)

}
