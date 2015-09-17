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

	// create and assign a new watcher
	watcher, err := util.Watch()
	if err != nil {
		util.LogFatal("[commands/watch] fsnotify.NewWatcher() failed", err)
	}
	defer watcher.Close()

	fmt.Printf("\n%v", stylish.Bullet("Watching for chages at '%s'", config.CWDir))
	fmt.Printf("%v\n", stylish.SubBullet("(Ctrl + c to quit)"))

	for {
		select {

		// watch for events
		case event := <-watcher.Events:

			// don't care about chmod updates
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

				fmt.Printf("\n%v", stylish.Bullet("Watching for chages at '%s'", config.CWDir))
				fmt.Printf("%v\n", stylish.SubBullet("(Ctrl + c to quit)"))
			}

			// watch for errors
		case err := <-watcher.Errors:
			fmt.Println("WATCH ERROR!", err)
		}
	}
}
