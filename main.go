// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package main

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/commands"
	"github.com/pagodabox/nanobox-cli/config"
)

// init
func init() {

	// idempotent install
	config.Console.Debug("Verifying install...")
	if err := install(); err != nil {
		config.Console.Fatal("Failed to install! Exiting...", err)
		os.Exit(1)
	}

	// create a logger
	config.Console.Debug("Creating logger...")
	logger, err := lumber.NewFileLogger(config.LogFile, config.LogLevel, lumber.ROTATE, 100, 1, 100)
	if err != nil {
		config.Console.Fatal("Failed to create logger! Exiting...", err)
		os.Exit(1)
	}

	// set the logger
	config.Log = logger

	// set the api to debug mode
	if config.LogLevel == lumber.DEBUG {
		api.Debug = true
	}
}

// main
func main() {

	// check for updates
	// checkUpdate()

	// run the CLI
	run()
}

// run attempts to run a CLI command. If no flags are passed (only the program
// is run) it will default to printing the CLI help text. It takes a help flag
// for printing the CLI help text. It takes a version flag for displaying the
// current version. It takes an app flag to indicate which app to run the command
// on (otherwise it wll attempt to find an app associated with the current directory).
// It also takes a debug flag (which must be passed last), that will display all
// request/response output for any API call the CLI makes.
func run() {

	// command line args w/o program
	args := os.Args[1:]

	// if only program is run, print help by default
	if len(args) <= 0 {
		help()
	}

	// parse command line args; it's safe to assume that args[0] is the command we
	// want to run, or one of our 'shortcut' flags that we'll catch before trying
	// to run the command.
	command := args[0]

	// check for 'global' commands
	switch command {

	// check for help shortcuts
	case "-h", "--help", "help":
		help()

	// check for version shortcuts
	case "-v", "--version", "version":
		fmt.Printf("Version: %v\n", config.Version.String())

	// check for the update command and update
	case "-u", "--update", "update":
		update()

	// we didn't find a 'shortcut' flag, so we'll continue parsing the remaining
	// args looking for a command to run.
	default:

		// if we find a valid command we run it
		if val, ok := commands.Commands[command]; ok {

			// args[1:] will be our remaining subcommand or flags after the intial command.
			// This value could also be 0 if running an alias command.
			opts := args[1:]

			//
			if len(opts) >= 1 {
				switch opts[0] {

				// Check for help shortcuts on commands
				case "-h", "--help", "help":
					commands.Commands[command].Help()
					os.Exit(0)
				}
			}

			// if debug was passed we remove it from the list of options that get sent
			// to commands
			if args[len(args)-1] == "--debug" {
				opts = opts[:len(opts)-1]
			}

			// do a quick ping to make sure we can communicate properly with the API
			// if err := api.DoRawRequest(nil, "GET", "https://api.pagodabox.io/v1/ping", nil, nil); err != nil {
			// 	ui.LogFatal("[main] The CLI was unable to communicate with the API", err)
			// }

			// run the command
			val.Run(opts)

			// no valid command found
		} else {
			fmt.Printf("'%v' is not a valid command. Type 'nanobox' for available commands and usage.\n", command)
			os.Exit(1)
		}
	}
}

// help
func help() {
	cmd := commands.Commands["help"]
	cmd.Run(nil)
	os.Exit(0)
}
