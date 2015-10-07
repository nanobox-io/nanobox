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

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/notify"
	"github.com/nanobox-io/nanobox-cli/util/server"
	"github.com/nanobox-io/nanobox-golang-stylish"
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
	if err := notify.Watch(config.CWDir, server.NotifyRebuild); err != nil {
		fmt.Printf(stylish.ErrBullet("Unable to detect file changes - %v", err.Error()))
	}
}
