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

var timoutReader *TimeoutReader

type TimeoutReader struct {
	Files   chan string
	timeout time.Duration
	once    sync.Once
}

func (self *TimeoutReader) Read(buf []byte) (n int, err error) {
	n = 0
	err = nil
	select {
	case file := <-self.Files:
		// if i recieve a file on the channel let it be read
		n = copy(buf, file+"\n")
		if n < len(file+"\n") {
			err = fmt.Errorf("Filename not coppied")
		}
		return
	case <-time.After(self.timeout):
		// if the timeout happens close the connection and EOF
		timoutReader = nil
		self.once.Do(func() {
			close(self.Files)
		})
		return 0, io.EOF
	}
	return
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
	// if there is no timeout reader or open request create one
	if timoutReader == nil {
		// create a new timeout reader
		timoutReader = &TimeoutReader{
			Files:   make(chan string),
			timeout: 10 * time.Second,
		}

		go func() {
			if _, err := Post("/file-changes", "text/plain", timoutReader); err != nil {
				config.Error("file changes error", err.Error())
			}
		}()
	}

	// get the name from the event and put the full path on it
	name := strings.Replace(event.Name, config.CWDir, "", -1)

	// if it is not just a chmod send the file
	if event.Op != fsnotify.Chmod {
		timoutReader.Files <- name
	}

	return nil
}
