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
	"github.com/nanobox-io/nanobox-cli/util/vagrant"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var reloadCmd = &cobra.Command{
	Hidden: true,

	Use:   "reload",
	Short: "Reloads the nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    reload,
}

// reload runs 'vagrant reload --provision'
func reload(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	fmt.Printf(stylish.Bullet("Reloading nanobox..."))
	if err := vagrant.Reload(); err != nil {
		config.Fatal("[commands/reload] failed - ", err.Error())
	}
}
