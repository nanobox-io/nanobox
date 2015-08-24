// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package main

import (
	// "fmt"
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

	// do a quick ping to make sure we can communicate properly with the API
	// if err := api.DoRawRequest(nil, "GET", "https://api.pagodabox.io/v1/ping", nil, nil); err != nil {
	// 	ui.LogFatal("[main] The CLI was unable to communicate with the API", err)
	// }

	//
	commands.NanoboxCmd.Execute()
}
