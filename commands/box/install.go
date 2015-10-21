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
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,

	Run: Install,
}

// Install
func Install(ccmd *cobra.Command, args []string) {
	if err := checkInstall(); err != nil {
		Config.Fatal("[commands/boxInstall] failed - ", err.Error())
	}
}

func checkInstall() (err error) {
	// install the nanobox vagrant image only if it isn't already available
	if !Vagrant.HaveImage() {
		fmt.Printf(stylish.Bullet("Installing nanobox image..."))

		// install the nanobox vagrant image
		err = Vagrant.Install()
	}
	return
}
