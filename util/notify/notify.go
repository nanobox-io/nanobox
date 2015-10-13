// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package notify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/go-fsnotify/fsnotify"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var (
	watcher    *fsnotify.Watcher
	ignoreDirs = []string{}
)

//
func Watch(path string, handle func(e *fsnotify.Event) error) error {

	var err error

	//
	setFileLimit()

	// get a list of directories that should not be watched; this is done because
	// there is a limit to how many files can be watched at a time, so folders like
	// node_modules, bower_components, vendor, etc...
	if err = getIgnoreDirs(); err != nil {
		return err
	}

	// create a new file watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		if _, ok := err.(syscall.Errno); ok {
			fmt.Printf(`
! WARNING !
Failed to watch files, max file descriptor limit reached. Nanobox will not
be able to propagate filesystem events to the virtual machine. Consider
increasing your max file descriptor limit to re-enable this functionality.
`)
		}

		return err
	}

	//
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	// if the file is a directory, recursively add each subsequent directory to
	// the watch list
	if fi.Mode().IsDir() {
		if err = filepath.Walk(path, watchDir); err != nil {
			return err
		}

		// else, if the file is just a file, add only it to the watch list
	} else {
		if err = watcher.Add(path); err != nil {
			return err
		}
	}

	// watch for interrupts
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	// watch for file events (blocking)
	for {

		select {

		// handle any file events by calling the handler function
		case event := <-watcher.Events:

			switch event.Op {

			// the watcher needs to watch itself to see if any directories are added
			// to then add them to the list of watched files
			case fsnotify.Create:
				fi, err := os.Stat(event.Name)
				if err := watchDir(event.Name, fi, err); err != nil {
					fmt.Printf(stylish.ErrBullet("Unable to watch files - '%v'", err))
				}

			// the watcher needs to watch itself to see if any directories are removed
			// to then remove them from the list of watched files (this may need to happen
			// recursively)
			case fsnotify.Remove:
				if err = watcher.Remove(event.Name); err != nil {
					return err
				}
			}

			// call the handler for each even fired
			if err := handle(&event); err != nil {
				return err
			}

			// handle any errors by calling the handler function
		case err := <-watcher.Errors:
			fmt.Printf(stylish.ErrBullet("Unable to watch files - '%v'", err))

			// listen for any signals and retun execution back to the CLI to finish
			// w/e it might need to finish
		case <-exit:
			return nil
		}
	}
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() && !isIgnoreDir(fi.Name()) {
		if err = watcher.Add(path); err != nil {
			return err
		}
	}

	return nil
}

// isIgnoreDir
func isIgnoreDir(name string) bool {
	for _, dir := range ignoreDirs {
		if dir == name {
			return true
		}
	}
	return false
}

// getIgnoreDirs
func getIgnoreDirs() error {
	res, err := http.Get(fmt.Sprintf("%s/libdirs", config.ServerURL))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &ignoreDirs)
}
