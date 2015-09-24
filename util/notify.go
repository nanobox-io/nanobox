package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"github.com/go-fsnotify/fsnotify"

	"github.com/pagodabox/nanobox-cli/config"
)

var watcher *fsnotify.Watcher
var ignoreDirs = []string{}

//
func Watch(path string, fn func(e *fsnotify.Event, err error)) error {
	setFileLimit()
	// get a list of directories that should not be watched; this is done because
	// there is a limit to how many files can be watched at a time, so folders like
	// node_modules, bower_components, vendor, etc...
	if err := getIgnoreDirs(); err != nil {
		return err
	}

	// create a new file watcher
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// this is dumb
	watcher = w

	// recursivly walk directories, beginning at 'root' (path), adding watchers
	if err := filepath.Walk(path, watchDir); err != nil {
		return err
	}

	// watch for file events (blocking)
	for {
		select {

		// handle any file events by calling the handler function
		case event := <-watcher.Events:

			// the watcher needs to watch itself to see if any directories are added to then
			// add them to the list of watched files;
			if event.Op == fsnotify.Create {
				fi, err := os.Stat(event.Name)
				if err := watchDir(event.Name, fi, err); err != nil {
					fn(&event, err)
				}
			}

			// call the handler for each even fired
			fn(&event, nil)

			// handle any errors by calling the handler function
		case err := <-watcher.Errors:
			fn(nil, err)
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
	res, err := http.Get(fmt.Sprintf("http://%s/libdirs", config.ServerURI))
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

// set my rlimit to the maximum rlimit
func setFileLimit() {
	rlm := &syscall.Rlimit{}
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, rlm)
	rlm.Cur = rlm.Max
	syscall.Setrlimit(syscall.RLIMIT_NOFILE, rlm)
}
