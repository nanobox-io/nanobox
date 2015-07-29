// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"

	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// ConsoleCommand satisfies the Command interface
type ConsoleCommand struct{}

// Help
func (c *ConsoleCommand) Help() {
	ui.CPrint(`
Description:
  Drops you into bash inside of the nanobox vm docker

Usage:
  nanobox console
  `)
}

// Run
func (c *ConsoleCommand) Run(opts []string) {
	fmt.Printf(stylish.Bullet("Opening nanobox console"))
	exec := ExecCommand{ console: true }
	exec.Run(opts)
}
