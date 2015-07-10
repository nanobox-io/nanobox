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
	"github.com/pagodabox/nanobox-golang-stylish"
)

// Install
func install() error {

	// check for a ~/.nanobox dir and create one if it's not found
	if di, _ := os.Stat(config.NanoDir); di == nil {

		//
		fmt.Printf(stylish.Bullet("Creating " + config.NanoDir + " directory"))

		if err := os.Mkdir(config.NanoDir, 0755); err != nil {
			return err
		}
	}

	// check for a ~/.nanobox/apps dir and create one if it's not found
	if di, _ := os.Stat(config.AppsDir); di == nil {

		//
		fmt.Printf(stylish.Bullet("Creating " + config.AppsDir + " directory"))

		if err := os.Mkdir(config.AppsDir, 0755); err != nil {
			return err
		}
	}

	// check for a ~/.nanobox/.auth file and create one if it's not found
	if fi, _ := os.Stat(config.AuthFile); fi == nil {

		//
		fmt.Printf(stylish.Bullet("Creating " + config.AuthFile + " directory"))

		if _, err := os.Create(config.AuthFile); err != nil {
			return err
		}
	}

	// check for a ~/.nanobox/.update file and create one if it's not found
	if fi, _ := os.Stat(config.UpdateFile); fi == nil {

		//
		fmt.Printf(stylish.Bullet("Creating " + config.UpdateFile + " directory"))

		if _, err := os.Create(config.UpdateFile); err != nil {
			return err
		}
	}

	// create/override nanobox.log file (this should be created by the logger, but
	// just incase it's not we'll leave this here)
	// if fi, _ := os.Stat(config.LogFile); fi == nil {

	//  //
	//  fmt.Printf(stylish.Bullet("Creating " + config.LogFile + " directory"))

	//  if _, err := os.Create(config.LogFile); err != nil {
	//    return err
	//  }
	// }

	return nil
}
