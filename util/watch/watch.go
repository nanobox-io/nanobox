package watch

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
)

var ignoreFile = []string{".git", ".hg", ".svn", ".bzr"}
var changeList = []string{}

// the watch package watches a folder and all its sub folders
// in doing so it may run into open file errors or things of that nature
// if it does, it will automatically fall back to a slower but still
// useful watching mechanism that looks at files and gets a hash of the content
// then when

type event struct {
	file  string
	error error
}

type Watcher interface {
	watch() error
	eventChan() chan event
	close() error
}

// watch a directory and report changes to nanobox
func Watch(container, path string) error {
	populateIgnore(path)
	// try watching with the fast one
	watcher := newNotifyWatcher(path)
	err := watcher.watch()
	if err != nil {
		// if it fails display a message and try the slow one
		fmt.Println("fast watcher broke.. falling back to slow")
		lumber.Info("Error occured in fast notify watcher: %s", err.Error())

		watcher.close()
		watcher = newCrawlWatcher(path)
		err := watcher.watch()
		if err != nil {
			// neither of the watchers are working
			return err
		}
	}
	defer watcher.close()

	go batchPublish(container)

	// catch a kill signal
	for e := range watcher.eventChan() {
		containerFile := filepath.ToSlash(filepath.Join("/app", strings.Replace(e.file, config.LocalDir(), "", 1)))
		changeList = append(changeList, containerFile)
	}

	// report any errors from either
	fmt.Println("done??")
	return nil
}

func batchPublish(container string) {
	for {
		<-time.After(time.Second)
		if len(changeList) > 0 {
			util.DockerExec(container, "touch", changeList, nil)
			fmt.Println("updates!", changeList)
			changeList = []string{}
		}
	}
}

// populate the ignore file from the nanoignore
func populateIgnore(path string) {
	b, err := ioutil.ReadFile(filepath.ToSlash(filepath.Join(path, ".nanoignore")))
	if err != nil {
		return
	}

	stringFields := strings.Fields(string(b))
	for _, str := range stringFields {
		ignoreFile = append(ignoreFile, str)
	}
}
