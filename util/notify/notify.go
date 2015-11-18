//
package notify

import (
	"encoding/json"
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"github.com/nanobox-io/nanobox/config"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	watcher    *fsnotify.Watcher
	ignoreDirs = []string{}
)

// Watch
func Watch(path string, handle func(e *fsnotify.Event) error) error {

	var err error

	//
	setFileLimit()

	// get a list of directories that should not be watched; this is done because
	// there is a limit to how many files can be watched at a time, so folders like
	// node_modules, bower_components, vendor, etc...; if this fails we probably
	// want to return an error because we would be detremental to watch those files
	if err := getIgnoreDirs(); err != nil {
		return fmt.Errorf("Unable to read ignore dirs - %v", err.Error())
	}

	// add source control files to be ignored (git, mercuriel, svn)
	ignoreDirs = append(ignoreDirs, ".git", ".hg", "trunk")

	// create a new file watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		if _, ok := err.(syscall.Errno); ok {
			return fmt.Errorf(`
! WARNING !
Failed to watch files, max file descriptor limit reached. Nanobox will not
be able to propagate filesystem events to the virtual machine. Consider
increasing your max file descriptor limit to re-enable this functionality.
`)
		}

		// return error here because w/o a watcher we really cant do anything
		return fmt.Errorf("Failed to create watcher - %v", err.Error())
	}

	// return this err because that means the path to the file they are trying to
	// watch doesn't exist
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	switch {

	// if the file is a directory, recursively add each subsequent directory to
	// the watch list; fsnotify will watch all files in a directory
	case fi.Mode().IsDir():
		if err = filepath.Walk(path, watchDir); err != nil {
			return err
		}

	// if the file is just a file, add only it to the watch list
	case fi.Mode().IsRegular():
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

			//
			fi, err := os.Stat(event.Name)

			switch event.Op {

			// the watcher needs to watch itself to see if any files are added to then
			// add them to the list of watched files
			case fsnotify.Create:

				// ensure that the file still exists before trying to watch it; ran into
				// a case with VIM where some tmp file (.swpx) was create and removed in
				// the same instant causing the watch to panic
				if err == nil && fi != nil && fi.Mode().IsDir() {

					// just ignore errors here since there isn't really anything that can
					// be done about it
					watchDir(event.Name, fi, err)
				}

			// the watcher needs to watch itself to see if any directories are removed
			// to then remove them from the list of watched files
			case fsnotify.Remove:

				// ensure thath the file is still available to be removed before attempting
				// to remove it; the main reason for manually removing files is to help
				// spare the ulimit; just log errors here, dont disrupt workflow.
				if err == nil {
					if err := watcher.Remove(event.Name); err != nil {
						config.Info(fmt.Sprintf("Unable to watch file %v - %v", event.Name, err.Error()))
					}
				}
			}

			// call the handler for each even fired; just log errors here, dont disrupt
			// workflow.
			if err := handle(&event); err != nil {
				config.Info(fmt.Sprintf("Unable to stop watching file %v - %v", event.Name, err.Error()))
			}

		// handle any errors by calling the handler function
		case <-watcher.Errors:
			// do something with watch errors?

			// listen for any signals and retun execution back to the CLI to finish
			// w/e it might need to finish
		case <-exit:
			return nil
		}
	}
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// don't walk any directory that is an ignore dir
	if isIgnoreDir(fi.Name()) {
		return filepath.SkipDir
	}

	// recursively add watchers to directores only (fsnotify will watch all files
	// in an added directory). Also, dont watch any files/dirs on the ignore list;
	// just log errors here, dont disrupt workflow.
	if fi.Mode().IsDir() {
		if err = watcher.Add(path); err != nil {
			config.Info(fmt.Sprintf("Unable to watch dir %v - %v", path, err.Error()))
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
