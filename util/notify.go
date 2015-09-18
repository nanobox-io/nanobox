package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"

	"github.com/pagodabox/nanobox-cli/config"
)

var watcher *fsnotify.Watcher
var libDirs = []string{}

func WatchCWD(fn func(e *fsnotify.Event, err error)) {
	watcher, err := watch()
	if err != nil {
		fn(nil, err)
		return
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Op == fsnotify.Create {
				file, err := os.Open(event.Name)
				if err == nil {
					fi, err := file.Stat()
					if err == nil && fi.Mode().IsDir() {
						watcher.Add(event.Name)
					}
				}
			}
			fn(&event, nil)
		case err := <-watcher.Errors:
			fn(nil, err)
		}
	}
}

func watch() (*fsnotify.Watcher, error) {
	SetLibDirs()

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := filepath.Walk(config.CWDir, watchDir); err != nil {
		return nil, err
	}
	return watcher, nil
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() && !isLibDir(fi.Name()) {
		if err = watcher.Add(path); err != nil {
			return err
		}
	}

	return nil
}

func isLibDir(name string) bool {
	for _, libDir := range libDirs {
		if libDir == name {
			return true
		}
	}
	return false
}

func SetLibDirs() {
	resp, err := http.Get(fmt.Sprintf("http://%s/libdirs", config.ServerURI))
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	json.Unmarshal(b, &libDirs)
	return
}
