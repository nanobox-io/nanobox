// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"github.com/pagodabox/nanobox-cli/ui"
)

// HelpCommand satisfies the Command interface for obtaining user info
type HelpCommand struct{}

// Help prints detailed help text for the user command
func (c *HelpCommand) Help() {
	ui.CPrint(`
Description:
  Prints help text for entire CLI

Usage:
  pagoda
  pagoda help
  pagoda -h
  pagoda --help

  ex. pagoda update
  `)
}

// Run prints out the help text for the entire CLI
func (c *HelpCommand) Run(opts []string) {
	ui.CPrint(`

                                     ***
                                  *********
                             *******************
                         ***************************
                             *******************
                         ...      *********      ...
                             ...     ***     ...
                         +++      ...   ...      +++
                             +++     ...     +++
                         \\\      +++   +++      ///
                             \\\     +++     ///
                                  \\     //
                                     \//

                      _  _ ____ _  _ ____ ___  ____ _  _
                      |\ | |__| |\ | |  | |__) |  |  \/
                      | \| |  | | \| |__| |__) |__| _/\_



Description:
  Welcome to the nanobox CLI! This will be your primary tool when working with
  nanobox. If you encounter any issues or have any suggestions, [green]find us on
  IRC (freenode) at #nanobox[reset]. Our engineers are available between 8 - 5pm MST.

  All commands have a short [-*] and a verbose [--*] option when passing flags.

  You can pass -h, --help, or help to any command to receive detailed information
  about that command.

  You can pass --debug at the end of any command to see all request/response
  output when making API calls.

Usage:
  pagoda (<COMMAND>:<ACTION> OR <ALIAS>) [GLOBAL FLAG] <POSITIONAL> [SUB FLAGS] [--debug]

Options:
  -h, --help, help
    Run anytime to receive detailed information about a command.

  -v, --version, version
    Run anytime to see the current version of the CLI.

  --debug
    Shows all API request/response output. [red]MUST APPEAR LAST[reset]

Available Commands:

  create      : Runs an 'init' and starts a nanobox VM
  deploy      : Deploys to your nanobox VM
  destroy     : Destroys the current nanobox VM
  halt        : Halts the current nanobox VM
  help        : Display this help
  init        : Creates a nanobox flavored Vagrantfile
  log         : Show/Stream nanobox logs
  new         : Create a new nanobox developer project
  publish     : Publish your nanobox live
  reload      :
  resume      : Resumes a halted/suspened nanobox VM
  status      : Display all current nanobox VM's
  suspend     : Suspends the current nanobox VM
  up          : Runs a 'create' and a 'deploy'
  update      : Updates this CLI to the newest version
  `)
}
