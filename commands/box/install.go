// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package box

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/vagrant"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,

	Run: Install,
}

// Install
func Install(ccmd *cobra.Command, args []string) {

	// install the nanobox vagrant image only if it isn't already available
	if !vagrant.HaveImage() {
		fmt.Printf(stylish.Bullet("Installing nanobox image..."))

		// install the nanobox vagrant image
		if err := vagrant.Install(); err != nil {
			config.Fatal("[commands/boxInstall] failed - ", err.Error())
		}
	}
}
