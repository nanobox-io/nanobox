package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// AppCreateCommand satisfies the Command interface for creating an app
type AppCreateCommand struct{}

// Help prints detailed help text for the app create command
func (c *AppCreateCommand) Help() {
	ui.CPrintln(`
Description:
  Creates a new application on Nanobox.

  If [app-name] is not specified, a name will be generated for you.

  If [--tinker] is not specified, the app will be created as 'production'.

Usage:
  pagoda create [-a app-name] [-t]
  pagoda app:create [-a app-name] [-t]

  ex. pagoda create -a app-name -t

Options:
  -a, --app [app-name]
    The name of your app

  -t, --tinker
    By default, app's are created as 'production'. Passing
    the tinker flag will create your new app as a 'tinker' app.

    Tinker : Free app with limited functionality - 3 per user.
    Production : Paid app with full functionality.
  `)
}

// Run attempts to create a new app on Nanobox. It can take an app-name flag
// for naming the app, and a tinker flag for designating the type of app to create.
// If successful, it attempts to add a new remote, then prints instructions on
// pushing code to pagodabox
func (c *AppCreateCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fTinker bool
	flags.BoolVar(&fTinker, "t", false, "")
	flags.BoolVar(&fTinker, "tinker", false, "")

	if err := flags.Parse(opts); err != nil {
		fmt.Println("Unable to parse flags. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda app:create", err)
	}

	//
	appCreateOptions := &nanoAPI.AppCreateOptions{Name: fApp, Free: fTinker}

	// create app
	app, err := api.CreateApp(appCreateOptions)
	if err != nil {
		_, err, msg := helpers.HandleAPIError(err)
		fmt.Printf("Oops! We could not create your app: %v - %v", err, msg)
		os.Exit(1)
	}

	ui.CPrintln(`
[green]` + app.Name + `[reset] has been created!

Attempting to add git remote 'pagoda'...`)

	remote := "pagoda"
	remotePath := "git@git.pagodabox.io:apps/" + app.Name + ".git"

	// attempt to find a git config file, if none is found (err) inform and instruct
	// on how to init the current dir
	gitConfigFile, path, err := helpers.FindGitConfigFile()
	if err != nil {
		fmt.Printf(`
Unable to add remote. It appears you are not in a git repo. To initialize the
current working directory as a git repo and add a pagoda remote, run:
  git init
  git remote add pagoda git@git.pagodabox.io:apps/` + app.Name + `.git
    `)
		os.Exit(1)
	}

	// look for a pagoda remote in the .git/config file, if one is found prompt for
	// a new remote name
	if _, ok := gitConfigFile.Get(`remote "pagoda"`, "url"); ok {
		remote = ui.Prompt("Looks like you already have a(n) 'pagoda' remote. Please specify a new remote name: ")
	}

	// add a new remote
	cmd := exec.Command("git", "remote", "add", remote, remotePath)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		fmt.Printf("Unable to add remote at this time: %v \n", err)
		os.Exit(1)
	}

	ui.CPrintln(`
New remote added:
  [green]` + remote + ` ` + remotePath + `[reset]

To push code to your new app run:
  git add .
  git commit -am "Pushing to Nanobox\!"
  git push ` + remote + ` master
  `)
}
