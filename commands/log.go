// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	// "fmt"

	// api "github.com/pagodabox/nanobox-api-client"
	// "github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

// LogCommand satisfies the Command interface for listing a user's apps
type LogCommand struct{}

// Help prints detailed help text for the app list command
func (c *LogCommand) Help() {
	ui.CPrint(`
Description:
  Show/Stream nanobox logs

  If [count] is not provided, will show the last 100 lines of the log.

  If [live] is not provided, will default to showing the last 100 lines.

Usage:
  nanobox log [-c count] [-l]

  ex. nanobox log -c 100 -l

Options:
  -c, --count [count]
    The number of lines of the log you wish to view.

  -l, --live
    Enable live stream
  `)
}

// Run displays select information about all of a user's apps
func (c *LogCommand) Run(opts []string) {

}
