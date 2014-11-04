package commands

import (
	"flag"
	"fmt"
	"os"

	pagodaAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// AppDestroyCommand satisfies the Command interface for destroying an app
type AppDestroyCommand struct{}

// Help prints detailed help text for the app destroy command
func (c *AppDestroyCommand) Help() {
	ui.CPrintln(`
Description:
  Destroys an application on Pagoda Box. [red]THIS ACTION CANNOT BE UNDONE![reset]

  If [app-name] is not specified, will prompt for one. Also, before an app is
  destroyed, will prompt for [app-name] to confirm action.

Usage:
  pagoda destroy [-a app-name]
  pagoda app:destroy [-a app-name]

  ex. pagoda destroy -a app-name

Options:
  -a, --app [app-name]
    The name of the app you want to destroy.

  -f, --force
    A force delete [red]skips confirmation... use responsibly[reset]!
  `)
}

// Run attempts to destroy an app on Pagoda Box. It can take a force flag that will
// skip the confirmation process, other wise will ask for confirmation by retyping
// the name of the app to be destroyed
func (c *AppDestroyCommand) Run(fApp string, opts []string, api *pagodaAPI.Client) {

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fForce bool
	flags.BoolVar(&fForce, "f", false, "")
	flags.BoolVar(&fForce, "force", false, "")

	if err := flags.Parse(opts); err != nil {
		fmt.Println("Unable to parse flags. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda app:destroy", err)
	}

	// if no app flag was passed, prompt for which app they want to destroy
	if fApp == "" {
		fApp = ui.Prompt("Which app would you like to delete: ")
	}

	//
	if fForce {
		c.forceDeleteApp(api, fApp)

		//
	} else {
		c.safeDeleteApp(api, fApp)
	}

}

// forceDeleteApp skips the confirmation process and deletes the app
func (c *AppDestroyCommand) forceDeleteApp(api *pagodaAPI.Client, app string) {
	c.deleteApp(api, app)
}

// safeDeleteapp asks the user to confirm the name of the app before deletion
func (c *AppDestroyCommand) safeDeleteApp(api *pagodaAPI.Client, app string) {

	ui.CPrintln(`
All code, data, network storage, ect. will be deleted forever.
Are you sure you want to destroy [red]` + app + `[reset]?
  `)

	response := ui.CPrompt("To confirm, type the name of the app you'd like to destroy ([red]this CANNOT be undone![reset]): ")

	//
	if response != app {
		ui.CPrintln("You typed [blue]'" + response + "[reset]' which [yellow]does NOT match[reset] '[green]" + app + "[reset]', your app [yellow]has NOT been destroyed[reset].\n")
		os.Exit(1)
	}

	//
	c.deleteApp(api, app)
}

// deleteApp deletes the specified app
func (c *AppDestroyCommand) deleteApp(api *pagodaAPI.Client, app string) {

	err := api.DeleteApp(app)
	if err != nil {
		_, err, msg := helpers.HandleAPIError(err)
		fmt.Printf("Oops! We could not delete '%s': %s - %s", app, err, msg)
		os.Exit(1)
	}

	ui.CPrintln("[red]" + app + "[reset] destroyed forever... how sad :(")
}
