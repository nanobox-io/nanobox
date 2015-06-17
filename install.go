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

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

// Install
func install() error {

	//
	config.Console.Info("[install] Verifying install...")
	config.Console.Info("[install] Current version %v", config.Version.String())

	// check for a ~/.nanobox dir and create one if it's not found
	if di, _ := os.Stat(config.NanoDir); di == nil {

		//
		config.Console.Info("[install] Creating %v directory", config.NanoDir)

		if err := os.Mkdir(config.NanoDir, 0755); err != nil {
			return err
		}
	}

	// check for a ~/.nanobox/apps dir and create one if it's not found
	if di, _ := os.Stat(config.AppsDir); di == nil {

		//
		config.Console.Info("[install] Creating %v directory", config.AppsDir)

		if err := os.Mkdir(config.AppsDir, 0755); err != nil {
			return err
		}
	}

	// check for a ~/.nanobox/.auth file and create one if it's not found
	if fi, _ := os.Stat(config.AuthFile); fi == nil {

		//
		config.Console.Info("[install] Creating %v file", config.AuthFile)

		if _, err := os.Create(config.AuthFile); err != nil {
			return err
		}
	}

	// check for a ~/.nanobox/.update file and create one if it's not found
	if fi, _ := os.Stat(config.UpdateFile); fi == nil {

		//
		config.Console.Info("[install] Creating %v file", config.UpdateFile)

		if _, err := os.Create(config.UpdateFile); err != nil {
			return err
		}
	}

	// create/override nanobox.log file
	// if fi, _ := os.Stat(config.LogFile); fi == nil {

	//  //
	//  config.Console.Info("[install] Creating %v file", config.LogFile)

	//  if _, err := os.Create(config.LogFile); err != nil {
	//    return err
	//  }
	// }

	return nil
}

//
func uninstall(force bool) {

	//
	if force != true {

		response := ui.Prompt("Are you sure you want to uninstall the Pagoda Box CLI (y/N)? ")

		if response != "y" {
			fmt.Printf("'%v' - Pagoda Box CLI will not be uninstalled. Exiting...\n", response)
			os.Exit(0)
		}
	}

	fmt.Print("Uninstalling... ")

	//
	if err := os.RemoveAll(config.NanoDir); err != nil {
		ui.LogFatal("[install] os.Remove() failed", err)
	}

	ui.CPrint("[green]success[reset]")
}
