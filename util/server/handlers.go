// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

// these are handlers that are passed into the util/notify/Watch command; they are
// called each time a file event happens
package server

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-fsnotify/fsnotify"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/server/mist"
)

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

	//
	name := strings.Replace(event.Name, config.CWDir, "", -1)

	if event.Op != fsnotify.Chmod {
		if _, err := Post("/file-changes?filename="+name, "text/plain", nil); err != nil {
			return err
		}
	}

	return nil
}
