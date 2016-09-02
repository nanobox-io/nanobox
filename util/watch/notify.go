package watch

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/jcelliott/lumber"
)

type notify struct {
	path    string
	events  chan event
	watcher *fsnotify.Watcher
}

func newNotifyWatcher(path string) Watcher {
	return &notify{
		path:   path,
		events: make(chan event, 10),
	}
}

// start the watching process and return an error if we cant watch all the files
func (n *notify) watch() (err error) {

	n.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = filepath.Walk(n.path, n.walkFunc)
	if err != nil {
		return err
	}

	go n.EventHandler()

	return
}

func (n *notify) eventChan() chan event {
	return n.events
}

// close the watcher
func (n *notify) close() error {
	return n.watcher.Close()
}

// add a file that is being walked to the watch system
func (n *notify) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	for _, ignoreName := range ignoreFile {
		if strings.HasSuffix(path, ignoreName) {
			return filepath.SkipDir
		}
	}

	return n.watcher.Add(path)
}

func (n *notify) EventHandler() {
	for {
		select {
		case e := <-n.watcher.Events:
			lumber.Debug("e: %+v", e)
			switch {
			case e.Op&fsnotify.Create == fsnotify.Create:
				// a new file/folder was created.. add it
				n.watcher.Add(e.Name)
				// send an event
				n.events <- event{file: e.Name}

			case e.Op&fsnotify.Write == fsnotify.Write:
				// a file was written to.. send the event
				n.events <- event{file: e.Name}

			case e.Op&fsnotify.Remove == fsnotify.Remove:
				// a file was removed. remove it
				n.watcher.Remove(e.Name)

			case e.Op&fsnotify.Rename == fsnotify.Rename:
				// remove from watcher because we no longer need to watch
				n.watcher.Remove(e.Name)

			case e.Op&fsnotify.Chmod == fsnotify.Chmod:
				// ignore anything that is just changing modes
				// mostlikely it was just a touch

			}
		case err := <-n.watcher.Errors:
			n.events <- event{file: "", error: err}
		}
	}

}
