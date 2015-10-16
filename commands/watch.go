// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"fmt"
<<<<<<< HEAD

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/notify"
	"github.com/nanobox-io/nanobox/util/server"
=======
>>>>>>> using interfaces for everything
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
)

//
var watchCmd = &cobra.Command{
	Hidden: true,

	Use:   "watch",
	Short: "",
	Long:  ``,

	Run: watch,
}

// watch
func watch(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Watching app files for changes"))

	// begin watching for file changes at cwd
	if err := Notify.Watch(config.CWDir, Server.NotifyRebuild); err != nil {
		fmt.Printf(err.Error())
	}
}
