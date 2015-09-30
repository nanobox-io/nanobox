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

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var watchCmd = &cobra.Command{
	Hidden: true,

	Use:   "watch",
	Short: "",
	Long:  ``,

	Run: nanoWatch,
}

// nanoWatch
func nanoWatch(ccmd *cobra.Command, args []string) {

	// indicate that files are being watched
	fWatch = true

	//
	fmt.Printf(stylish.Bullet("Watching app files for changes"))

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
				fmt.Printf(stylish.Bullet("Rebuilding environment"))
				nanoDeploy(nil, args)

				// ...otherwise just build
			} else {
				fmt.Printf(stylish.Bullet("Rebuilding code"))
				nanoBuild(nil, args)
			}
		}
	}); err != nil {

		//
		if _, ok := err.(syscall.Errno); ok {
			fmt.Printf(`
! WARNING !
Failed to watch files, max file descriptor limit reached. Nanobox will not
be able to propagate filesystem events to the virtual machine. Consider
increasing your max file descriptor limit to re-enable this functionality.
`)
		} else {
			fmt.Printf(stylish.ErrBullet("Unable to detect file changes (%v)", err.Error()))
		}

		os.Exit(1)
	}
}
