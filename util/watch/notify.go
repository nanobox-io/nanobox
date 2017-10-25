package watch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/jcelliott/lumber"
)

type event struct {
	file  string
	error error
	fsnotify.Event
}

type Watcher interface {
	watch() error
	eventChan() chan event
	close() error
}

type notify struct {
	events chan event // separate event channel so we don't send on all fsnotify.Events
	*fsnotify.Watcher
}

func newRecursiveWatcher(path string) (Watcher, error) {
	folders := subfolders(path)
	if len(folders) == 0 {
		return nil, fmt.Errorf("No folders to watch")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	notifyWatcher := &notify{Watcher: watcher, events: make(chan event)}

	for i := range folders {
		lumber.Info("Adding %s", folders[i])
		err = notifyWatcher.Add(folders[i])
		if err != nil {
			return nil, err
		}
	}
	return notifyWatcher, nil
}

func run(watcher *notify) {
	for {
		select {
		case evnt := <-watcher.Events:
			if shouldIgnoreFile(filepath.Base(evnt.Name)) {
				continue
			}

			if evnt.Op&fsnotify.Create == fsnotify.Create {
				fi, err := os.Stat(evnt.Name)
				if err != nil {
					// stat ./4913: no such file or directory
					lumber.Error("Failed to stat on event - %s", err.Error())
				} else if fi.IsDir() {
					lumber.Info("Detected dir creation: %s", evnt.Name) // todo: Debug doesn't work
					if !shouldIgnoreFile(filepath.Base(evnt.Name)) {
						err = watcher.Add(evnt.Name)
						if err != nil {
							lumber.Error("ERROR - %s", err.Error())
						}
					}
				} else {
					lumber.Info("Detected     creation: %s", evnt.Name) // todo: Debug doesn't work
					watcher.events <- event{file: evnt.Name}
				}
			}

			if evnt.Op&fsnotify.Write == fsnotify.Write {
				lumber.Info("Detected modification: %s", evnt.Name) // todo: Debug doesn't work
				watcher.events <- event{file: evnt.Name}
			}

		case err := <-watcher.Errors:
			lumber.Error("Watcher error encountered - %s", err.Error())
		}
	}
}

// subfolders returns a slice of subfolders (recursive), including the folder provided.
func subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			// skip folders that begin with a dot
			if shouldIgnoreFile(name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

// shouldIgnoreFile determines if a file should be ignored.
// Ignore files that start with `.` or `_` or end with `~`.
func shouldIgnoreFile(name string) bool {
	for i := range ignoreFile {
		if name == ignoreFile[i] {
			return true
		}
	}

	return strings.HasPrefix(name, ".") ||
		strings.HasPrefix(name, "_") ||
		strings.HasSuffix(name, ".swp") ||
		strings.HasSuffix(name, "~")
}

// start the watching process and return an error if we cant watch all the files
func (n *notify) watch() (err error) {
	go run(n)
	return nil
}

func (n *notify) eventChan() chan event {
	return n.events
}

// close the watcher
func (n *notify) close() error {
	return n.Close()
}
