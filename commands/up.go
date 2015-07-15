// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	// "fmt"

	// "github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

// UpCommand satisfies the Command interface
type UpCommand struct{}

// Help prints detailed help text for the app list command
func (c *UpCommand) Help() {
	ui.CPrint(`
Description:
  Runs a 'create' and a 'deploy'

Usage:
  nanobox up
  `)
}

// Run creates the specified virtual machine
func (c *UpCommand) Run(opts []string) {

	// run a create command to create a Vagrantfile and boot the VM...
	create := CreateCommand{}
	create.Run(opts)

	// ...create a deploy
	deploy := DeployCommand{}
	deploy.Run(opts)
}
