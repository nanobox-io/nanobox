// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// WatchCommand satisfies the Command interface
type WatchCommand struct {
	watcher *fsnotify.Watcher
}

// Help
func (c *WatchCommand) Help() {
	ui.CPrint(`
Description:
  Watches your app for file changes. When a file is changed a 'nanobox build' is
	automatically issued. If a Boxfile is modified a 'nanobox deploy' is issued.

Usage:
  nanobox watch
  `)
}

// Run
func (c *WatchCommand) Run(opts []string) {

	// create and assign a new watcher
	if watcher, err := fsnotify.NewWatcher(); err != nil {
		ui.LogFatal("[commands.watch] fsnotify.NewWatcher() failed", err)
	} else {
		c.watcher = watcher
	}
	defer c.watcher.Close()

	fmt.Printf("\n%v", stylish.Bullet(fmt.Sprintf("Watching for chages at '%v'", config.CWDir)))
	fmt.Printf("%v\n", stylish.SubBullet("(Ctrl + c to quit)"))

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(config.CWDir, c.watchDir); err != nil {
		ui.LogFatal("[commands.watch] filepath.Walk() failed", err)
	}

	for {
		select {

		// watch for events
		case event := <-c.watcher.Events:

			// don't care about chmod updates
			if event.Op != fsnotify.Chmod {

				// if the file changes is the Boxfile do a full deploy...
				if filepath.Base(event.Name) == "Boxfile" {
					fmt.Printf(stylish.Bullet("Issuing deploy"))
					deploy := DeployCommand{}
					deploy.Run(opts)

					// ...otherwise just build
				} else {
					fmt.Printf(stylish.Bullet("Issuing build"))
					build := BuildCommand{}
					build.Run(opts)
				}

				fmt.Printf("\n%v", stylish.Bullet(fmt.Sprintf("Watching for chages at '%v'", config.CWDir)))
				fmt.Printf("%v\n", stylish.SubBullet("(Ctrl + c to quit)"))
			}

			// watch for errors
		case err := <-c.watcher.Errors:
			fmt.Println("WATCH ERROR!", err)
		}
	}
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func (c *WatchCommand) watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		if err = c.watcher.Add(path); err != nil {
			return err
		}
	}

	return nil
}
