package commands

import (
	"fmt"
	"os"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// ServiceRebootCommand satisfies the Command interface for rebooting an app service
type ServiceRebootCommand struct{}

// Help prints detailed help text for the service reboot command
func (c *ServiceRebootCommand) Help() {
	ui.CPrintln(`
Description:
  Reboots a service in an app.

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

  The service-name/UID is [yellow]required[reset]. (Ex. web1)

Usage:
  pagoda service:reboot [-a app-name] service-name/UID

  ex. pagoda service:reboot -a app-name web1

Options:
  -a, --app [app-name]
    The name of the app
  `)
}

// Run attempts to reboot an app's service
func (c *ServiceRebootCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	var fService string

	// If there's no service, prompt for one
	if len(opts) <= 0 {
		fService = ui.Prompt("Which service would you like to reboot: ")

		// We should expect opts[0] to be the service.
	} else {
		fService = opts[0]
		opts = opts[1:]
	}

	service, err := helpers.GetServiceBySlug(fApp, fService, api)
	if err != nil {
		fmt.Printf("Oops! We could not find a '%v' on '%v'.\n", fService, fApp)
		os.Exit(1)
	}

	if _, err := api.RebootAppService(fApp, service.ID); err != nil {
		fmt.Println("There was a problem reboot %v. See ~/.pagodabox/log.txt for details", fService)
		ui.Error("pagoda service:reboot", err)
	}

	fmt.Printf("Rebooting %v's %v (%v)... Check your dashboard for transaction details and logs.\n", fApp, service.UID, service.Name)
}
