package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// AppOpenCommand satisfies the Command interface for opening a user's app
type AppOpenCommand struct{}

// Help prints detailed help text for the app open command
func (c *AppOpenCommand) Help() {
	ui.CPrintln(`
Description:
  Open an app in the default browser

Usage:
  pagoda open [-a app-name] [-p path]
  pagoda app:open [-a app-name] [-p path]

  ex. pagoda open -a app-name -p /cron-jobs

Options:
  -a, --app [app-name]
    The name of the app to open

  -p, --path [path]
    A specific path to open
  `)
}

// Run opens an app in the users default browser
func (c *AppOpenCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fPath string
	flags.StringVar(&fPath, "p", "", "")
	flags.StringVar(&fPath, "path", "", "")

	if err := flags.Parse(opts); err != nil {
		fmt.Println("Unable to parse flags. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda app:open", err)
	}

	//
	app, err := api.GetApp(fApp)
	if err != nil {
		fmt.Printf("Unable to find an app with the name '%v' \n", fApp)
		os.Exit(1)
	}

	fmt.Printf("Opening %v to %v", fApp, fPath)

	path := `https://dashboard.pagodabox.io/apps/` + app.ID + fPath

	// detect current operating system to determine which version of 'open' we can
	// call
	switch runtime.GOOS {
	case "linux":
		if err := exec.Command("xdg-open", path).Run(); err != nil {
			fmt.Println("There was a problem opening your app. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda app:open", err)
		}
	case "windows", "darwin":
		if err := exec.Command("open", path).Run(); err != nil {
			fmt.Println("There was a problem opening your app. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda app:open", err)
		}
	default:
		fmt.Println("There was a problem opening your app. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda app:open", err)
	}

	fmt.Println("")
}
