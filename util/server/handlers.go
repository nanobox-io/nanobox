// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package server

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-fsnotify/fsnotify"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/notify"
	"github.com/nanobox-io/nanobox-cli/util/server/mist"
)

// NotifyRebuild
func NotifyRebuild(event *fsnotify.Event) (err error) {
	if event.Op != fsnotify.Chmod {

		// pause logs
		config.Silent = true

		// the job thats going to be run; usually a build
		job := "build"

		switch filepath.Base(event.Name) {

		// run a build for any file changes
		default:
			fmt.Printf(`
++> BOXFILE CHANGED, CLOSING LOG STREAM FOR REBUILD ////////////////////////////
`)
		// if the file changes is the Boxfile, deploy
		case "Boxfile":
			fmt.Printf(`
++> SOURCE CODE CHANGED, CLOSING LOG STREAM FOR REBUILD ////////////////////////
`)
			job = "deploy"
		}

		done := make(chan struct{})

		// listen for status updates
		go func() {
			if err := mist.Listen([]string{"job", job}, mist.HandleDeployStream); err != nil {
				config.Fatal("[commands/nanoBuild] failed - ", err.Error())
			}
			close(done)
		}()

		// run 'job'
		switch job {
		case "build":
			err = Build("")
		case "deploy":
			err = Deploy("")
		}

		if err != nil {
			fmt.Printf(`
! AN ERROR PREVENTED NANOBOX FROM BUILDING YOUR ENVIRONMENT !
- View the output above to diagnose the source of the problem
- You can also retry with --verbose for more detailed output
`)
			//
			return notify.WatchError{"Watch failed"}
		}

		<-done

		fmt.Printf(`
--------------------------------------------------------------------------------
[√] APP SUCCESSFULLY REBUILT   ///   DEV URL : %v
--------------------------------------------------------------------------------

++> STREAMING LOGS (ctrl-c to exit) >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
`, config.Nanofile.Domain)

		// resume logs
		config.Silent = false
	}

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
