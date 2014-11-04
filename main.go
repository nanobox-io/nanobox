package main

import (
	"fmt"
	"os"
	"strings"

	pagodaAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

const Version = "0.0.1"

type (

	// CLU represents the Pagoda Box CLI. It has a version, a Pagoda Box API client
	// and a map of all the commands it responds to
	CLI struct {
		version   string
		apiClient *pagodaAPI.Client
		commands  map[string]Command
	}
)

// main creates a new CLI and then checks to see if authentication is needed. If
// no authentication is required it will attempt to run the provided command
func main() {

	//
	cli := &CLI{
		version:   Version,
		apiClient: pagodaAPI.NewClient(),
		commands:  Commands,
	}

	// see if there is a .pagodabox/token file. We wont handle an error here because
	// on a first time install we wouldn't expect to find one. Instead we'll handle
	// the file
	if helpers.NeedsAuth() {

		fmt.Println(`
Before you'll be able to use the Pagoda Box CLI on this machine. To continue,
please login to verify your account:
    `)

		// authenticate w/o username or password
		helpers.Authenticate("", "", cli.apiClient)

		fmt.Println("To begin using the Pagoda Box CLI type 'pagoda' to see a list of commands.")
		os.Exit(0)

		// we've already 'installed' the CLI so we can just create our client
	} else {

		// cli.apiClient.APIURL     = "https://dashboard.pagodabox.io"
		// cli.apiClient.APIVersion = "api"
		cli.apiClient.AuthToken = helpers.GetAuthToken()

		// run the CLI
		cli.run()
	}

}

// run attempts to run a CLI command. If no flags are passed (only the program
// is run) it will default to printing the CLI help text. It takes a help flag
// for printing the CLI help text. It takes a version flag for displaying the
// current version. It takes an app flag to indicate which app to run the command
// on (otherwise it wll attempt to find an app associated with the current directory).
// It also takes a debug flag (which must be passed last), that will display all
// request/response output for any API call the CLI makes.
func (cli *CLI) run() {

	// command line args w/o program
	args := os.Args[1:]

	// if only program is run, print help by default
	if len(args) <= 0 {
		cli.Help()

		// parse command line args
	} else {

		// it's safe to assume that args[0] is the command we want to run, or one of
		// our 'shortcut' flags that we'll catch before trying to run the command.
		command := args[0]

		// check for 'global' commands
		switch command {

		// Check for help shortcuts
		case "-h", "--help", "help":
			cli.Help()

		// Check for version shortcuts
		case "-v", "--version", "version":
			ui.CPrintln("[yellow]Version " + cli.Version() + "[reset]")

		// we didn't find a 'shortcut' flag, so we'll continue parsing the remaining
		// args looking for a command to run.
		default:

			// if we find a valid command we run it
			if val, ok := cli.commands[command]; ok {

				// args[1:] will be our remaining subcommand or flags after the intial command.
				// This value could also be 0 if running an alias command.
				opts := args[1:]

				// assume they wont be passing an app
				fApp := ""

				//
				if len(opts) >= 1 {
					switch opts[0] {

					// Check for help shortcuts
					case "-h", "--help", "help":
						cli.commands[command].Help()
						os.Exit(0)

					// Check for app flag, set fApp and strip out the flag and app
					case "-a", "--app":
						fApp = opts[1]
						opts = opts[2:]
					}
				}

				// before we run the command we'll check to see if debug mode needs to
				// be enabled. If so, enable it and strip off the flag.
				if args[len(args)-1] == "--debug" {
					cli.apiClient.Debug = true

					opts = opts[:len(opts)-1]
				}

				// do a quick ping to make sure we can communicate properly with the API
				_, err := cli.apiClient.GetUser("me")
				if err != nil {

					//
					if strings.Contains(err.Error(), "Invalid authentication token") {
						helpers.ReAuthenticate(cli.apiClient)

						//
					} else {
						fmt.Printf("The CLI was unable to communicate with the API: %s", err)
					}
				}

				// run the command
				val.Run(fApp, opts, cli.apiClient)

				// no valid command found
			} else {
				fmt.Printf("'%s' is not a valid command. Type 'pagoda' for available commands\n and usage.", command)
				os.Exit(1)
			}
		}
	}
}
