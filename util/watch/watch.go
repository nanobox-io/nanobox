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

// watch a directory and report changes to nanobox
func Watch(container, path string) error {
	populateIgnore(path)
	// ctimeCheck(container) // see comment in batchPublish()

	lumber.Info("watch: ignored files: %+v", ignoreFile)

	// try watching with the fast one
	watcher, err := newRecursiveWatcher("./")
	defer watcher.close()
	if err != nil {
		// if it fails display a message and try the slow one
		lumber.Info("Error occured in fast notify watcher: %s", err.Error())

		// print the warning
		// we added /r because this message often appears in a raw terminal which requires
		// carriage returns
		fmt.Printf("\n\r---------------------------------------------------------------------\n\r\n\r")
		fmt.Printf("An error occured in the fast notify watcher: '%s'\n\r\n\r", err.Error())
		fmt.Printf("Generally, having too low of a ulimit causes these issues.\n\r")
		fmt.Printf("Consider upping your ulimit to resolve (`ulimit -n 2048`).\n\r")
		fmt.Printf("Until then, we'll go ahead and rollover to a slower polling solution.\n\r")
		fmt.Printf("\n\r---------------------------------------------------------------------\n\r\n\r")

		watcher.close()
		watcher = newCrawlWatcher(path)
		err := watcher.watch()
		if err != nil {
			// neither of the watchers are working
			return err
		}
	} else {
		go run(watcher.(*notify))
	}

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

func uniq(list []string) []string {
	if len(list) == 0 {
		return []string{}
	}

	tmp := make(map[string]struct{})
	newlist := []string{}
	for i := range list {
		tmp[list[i]] = struct{}{}
	}
	for k, _ := range tmp {
		newlist = append(newlist, k)
	}
	return newlist
}

// publish in batches so to save clock cycles
func batchPublish(container string) {
	for {
		<-time.After(time.Second)
		changeList = uniq(changeList)
		if len(changeList) > 0 {
			lumber.Info("watcher: pushing: %+v", changeList)

			// As it turns out, the ctime behavior isn't fool-proof. After running
			// ctime as the default for a couple of weeks, reports have
			// confirmed that a majority of frameworks will not actually detect
			// a change if file attributes haven't been modified. We will need
			// to re-address this, either as a configurable option, or as a
			// custom patched kernel with a custom syscall that we can tie into.
			//
			// Until either of the aforementioned options are available, the
			// responsible decision here is to rollback to the previous point.
			//
			// if ctimeAvailable {
			// 	ctime(container, changeList)
			// } else {
			// 	touch(container, changeList)
			// }

			touch(container, changeList)

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
	out, err := util.DockerExec(container, "root", "touch", append([]string{"-c"}, changeList...), nil)
	if err != nil {
		fmt.Printf("TOUCH OUTPUT: %s\n", out)
		fmt.Printf("TOUCH  ERROR: %s\n", err.Error())
	}
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
