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
	// "os/exec"
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"

	// "github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	// "github.com/pagodabox/nanobox-golang-stylish"
)

// WatchCommand satisfies the Command interface
type WatchCommand struct{}

// Help prints detailed help text for the app list command
func (c *WatchCommand) Help() {}

var watcher *fsnotify.Watcher

// Run resumes the specified virtual machine
func (c *WatchCommand) Run(opts []string) {

	watcher, _ = fsnotify.NewWatcher()
	// if err != nil {
	// 	fmt.Println("BONK!", err)
	// }
	defer watcher.Close()

	watch := config.CWDir

	fmt.Println("WATCHING AT:", watch)

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(watch, watchDir); err != nil {
		ui.LogFatal("[commands.publish] filepath.Walk() failed", err)
	}

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {

			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)

				base := filepath.Base(event.Name)

				fmt.Println("BASE!", base)

				// if the name of the file is the type of  file I'm looking for (in this
				// case a 'Boxfile'), then do a full deploy otherwise do a build
				if base == "Boxfile" {
					fmt.Println("DO DEPLOY!")
				} else {
					fmt.Println("DO BUILD!")
				}

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR!", err)
			}
		}
	}()

	<-done

	fmt.Println("DONE!")
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	fmt.Println("FILE!", path)

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		fmt.Println("ADDING WATCHER TO: ", path)
		if err = watcher.Add(path); err != nil {
			return err
		}
	}

	return nil
}
