// these are handlers that are passed into the util/notify/Watch command; they are
// called each time a file event happens
package server

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-fsnotify/fsnotify"
	"github.com/nanobox-io/nanobox/util/server/mist"

	"github.com/nanobox-io/nanobox/config"
)

// A TimeoutReader reads from Files until timeout returning EOF
type TimeoutReader struct {
	Files   chan string
	timeout time.Duration
	once    sync.Once
}

var timoutReader *TimeoutReader

// Read
func (r *TimeoutReader) Read(p []byte) (n int, err error) {

	select {

	// if a file is received on the channel, read it
	case file := <-r.Files:
		n = copy(p, file+"\n")
		if n < len(file+"\n") {
			err = fmt.Errorf("Filename not coppied")
		}
		return

	// if no files are received before the timout send EOF to kill the connection
	case <-time.After(r.timeout):
		timoutReader = nil
		r.once.Do(func() {
			close(r.Files)
		})
		return 0, io.EOF
	}
}

// NotifyRebuild
func NotifyRebuild(event *fsnotify.Event) (err error) {

	// pause logs
	config.Silent = true

	//
	switch event.Op {

	// only care about create, write, remove, and rename events
	case fsnotify.Create, fsnotify.Write, fsnotify.Remove, fsnotify.Rename:

		//
		errch := make(chan error)

		switch filepath.Base(event.Name) {

		// run a build for any file changes
		default:
			fmt.Printf(`
++> SOURCE CODE CHANGED, CLOSING LOG STREAM FOR REBUILD ////////////////////////////
`)

			go func() {
				errch <- mist.Listen([]string{"job", "build"}, mist.BuildUpdates)
			}()

			//
			if err = Build(""); err != nil {
				return
			}

		// if the file changes is the Boxfile, deploy
		case "Boxfile":
			fmt.Printf(`
++> BOXFILE CHANGED, CLOSING LOG STREAM FOR REBUILD ////////////////////////
`)

			go func() {
				errch <- mist.Listen([]string{"job", "deploy"}, mist.DeployUpdates)
			}()

			//
			if err = Deploy(""); err != nil {
				return
			}
		}

		// wait for a status update (blocking)
		err = <-errch

		//
		if err != nil {
			return
		}

		fmt.Printf(`
--------------------------------------------------------------------------------
[√] APP SUCCESSFULLY REBUILT   ///   DEV URL : %v
--------------------------------------------------------------------------------

++> STREAMING LOGS (ctrl-c to exit) >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	`, config.Nanofile.Domain)

	}

	// resume logs
	config.Silent = false

	return
}

// NotifyServer
func NotifyServer(event *fsnotify.Event) error {

	// if there is no timeout reader create one and open a request; if there is no
	// timeout reader there wont be an open request, so checking for timeoutReader
	// is enough
	if timoutReader == nil {

		// create a new timeout reader
		timoutReader = &TimeoutReader{
			Files:   make(chan string),
			timeout: 10 * time.Second,
		}

		// launch a new request that is held open until EOF from the timeoutReader
		go func() {
			if _, err := Post("/file-changes", "text/plain", timoutReader); err != nil {
				config.Error("file changes error", err.Error())
			}
		}()
	}

	// strip the current working directory from the filepath
	path := strings.Replace(event.Name, config.CWDir, "", -1)

	// for any event other than Chmod, append the filepath to the list of files to
	// be read
	if event.Op != fsnotify.Chmod {
		timoutReader.Files <- path
	}

	return nil
}
