package watch

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/provider"
)

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
func Watch(path string) error {
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

	// catch a kill signal
	for e := range watcher.eventChan() {

		providerFile := filepath.ToSlash(filepath.Join(fmt.Sprintf("%s%s/code", provider.HostShareDir(), config.EnvID()), strings.Replace(e.file, config.LocalDir(), "", 1)))
		provider.Touch(providerFile)
	}

	// report any errors from either

	return nil
}
