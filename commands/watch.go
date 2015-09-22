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

	"github.com/go-fsnotify/fsnotify"
	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Resumes the halted/suspended nanobox VM",
	Long: `
Description:
  Watches your app for file changes. When a file is changed a 'nanobox build' is
  automatically issued. If a Boxfile is modified a 'nanobox deploy' is issued.`,

	Run: nanoWatch,
}

// nanoWatch
func nanoWatch(ccmd *cobra.Command, args []string) {

	fmt.Printf("\n%v", stylish.Bullet("Watching for chages at '%s' (Ctrl + c to quit)", config.CWDir))

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

			fmt.Printf("\n%v", stylish.Bullet("Watching for chages at '%s' (Ctrl + c to quit)", config.CWDir))
		}
	}); err != nil {
		fmt.Printf(stylish.ErrBullet("Unable to detect file changes (%v)", err.Error()))
		os.Exit(1)
	}
}
