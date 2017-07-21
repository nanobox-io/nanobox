package watch

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
)

var ignoreFile = []string{".git", ".hg", ".svn", ".bzr"}
var changeList = []string{}

var ctimeAvailable bool


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
	ctimeCheck(container)

	lumber.Debug("watch: ignored dirs: %+v", ignoreFile)
	// try watching with the fast one
	watcher := newNotifyWatcher(path)
	err := watcher.watch()
	if err != nil {
		// if it fails display a message and try the slow one
		lumber.Info("Error occured in fast notify watcher: %s", err.Error())

		// print the warning
		// we added /r because this message often appears in a raw terminal which requires
		// carrage returns
		fmt.Printf("\n\r-------------------------------------------------------\n\r\n\r")
		fmt.Printf("Uh oh, the live filesystem watcher has panic'ed.\n\r")
		fmt.Printf("We'll go ahead and rollover to a slower polling solution.\n\r")
		fmt.Printf("\n\r-------------------------------------------------------\n\r\n\r")

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
		efile := e.file
		if runtime.GOOS == "windows" {
			efile = strings.Replace(efile, "\\", "/", -1)
		}
		containerFile := filepath.ToSlash(filepath.Join("/app", strings.Replace(efile, config.LocalDir(), "", 1)))
		changeList = append(changeList, containerFile)
	}

	return nil
}

// publish in batches so to save clock cycles
func batchPublish(container string) {
	for {
		<-time.After(time.Second)
		if len(changeList) > 0 {
			lumber.Info("watcher: pushing: %+v", changeList)
			if ctimeAvailable {
				ctime(container, changeList)
			} else {
				touch(container, changeList)
			}

			changeList = []string{}
		}
	}
}

// check to see if ctime is installed on there docker image
func ctimeCheck(container string) {
	out, err := util.DockerExec(container, "root", "which", []string{"ctime"}, nil)
	if err == nil && strings.Contains(out, "ctime") {
		ctimeAvailable = true
	}
}

// the touch command used when ctime isnt available
func touch(container string, changeList []string) {
	util.DockerExec(container, "root", "touch", append([]string{"-c"}, changeList...), nil)
}

// the ctime command we will use
func ctime(container string, changeList []string) {
	util.DockerExec(container, "root", "ctime", changeList, nil)	
}

// populate the ignore file from the nanoignore
func populateIgnore(path string) {
	// add pieces from the env
	env, err := models.FindEnvByID(config.EnvID())
	box := boxfile.New([]byte(env.BuiltBoxfile))
	for _, libDir := range box.Node("run.config").StringSliceValue("cache_dirs") {
		ignoreFile = append(ignoreFile, libDir)
	}

	// add parts from the nanoignore
	b, err := ioutil.ReadFile(filepath.ToSlash(filepath.Join(path, ".nanoignore")))
	if err != nil {
		return
	}

	stringFields := strings.Fields(string(b))
	for _, str := range stringFields {
		ignoreFile = append(ignoreFile, str)
	}

}
