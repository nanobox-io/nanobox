// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/go-fsnotify/fsnotify"
	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var watchCmd = &cobra.Command{
	Hidden: true,

	Use:   "watch",
	Short: "",
	Long: `
Description:
  Watches your app for file changes. When a file is changed a 'nanobox build' is
  automatically issued. If a Boxfile is modified a 'nanobox deploy' is issued.`,

	Run: nanoWatch,
}

// nanoWatch
func nanoWatch(ccmd *cobra.Command, args []string) {

	// when watching files, we dont want to suspend the VM if a deploy fails
	fDebug = true

	//
	fmt.Printf("[√] Watching app files for changes\n")

	// begin watching for file changes at cwd
	if err := util.Watch(config.CWDir, func(event *fsnotify.Event, err error) {

		//
		if err != nil {
			fmt.Println(stylish.ErrBullet("Error detecting file change (%v)", err.Error()))
		}

		//
		if event.Op != fsnotify.Chmod {

			// if the file changes is the Boxfile do a full deploy...
			if filepath.Base(event.Name) == "Boxfile" {
				fmt.Printf(stylish.Bullet("Issuing deploy"))
				nanoDeploy(nil, args)

				// ...otherwise just build
			} else {
				fmt.Printf(stylish.Bullet("Issuing build"))
				nanoBuild(nil, args)
			}
		}
	}); err != nil {

		//
		if _, ok := err.(syscall.Errno); ok {
			fmt.Printf(stylish.ErrBullet("Insert error message for file error here"))
		}

		fmt.Printf(stylish.ErrBullet("Unable to detect file changes (%v)", err.Error()))
		os.Exit(1)
	}
}
