// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/ui"
)

// nanoHelp
func nanoHelp(ccmd *cobra.Command) error {
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
  IRC (freenode) at #nanobox[reset].

  All commands flags have a short [-*] and a verbose [--*] option.

  You can pass -h, --help, or help to any command to receive detailed information
  about that command.

  You can pass --debug [red]at the end[reset] of any command to see all request/response
  output when making API calls.

Usage:
  nanobox COMMAND [f, flag, -f, --flag] [--debug]

Options:
  h, -h, help, --help
    Run anytime to receive detailed information about a command.

  v, -v, version, --version
    Run anytime to see the current version of the CLI.

  --debug
    Shows all API request/response output [red](Must appear last)[reset].

Available Commands:

  bootstrap   : Runs an engine's bootstrap script - downloads code & launches VM
  build       : Rebuilds/Compiles your project
  console     : Opens an interactive terminal inside your apps context
  create      : Runs 'nanobox init' then boots the nanobox VM
  deploy      : Deploys to the nanobox VM
  destroy     : Destroys the nanobox VM
  exec        : Runs a command in your apps context
  fetch       : Fetches an engine from nanobox.io
  halt        : Halts the nanobox VM
  help        : Displays CLI help
  init        : Creates a nanobox flavored Vagrantfile
  log         : Shows/Streams nanobox logs
  new         : Generates a new engine
  publish     : Publishes an engine to nanobox.io
  reload      : Reloads the nanobox VM
  resume      : Resumes the halted/suspended nanobox VM
  status      : Displays all current nanobox VM's
  suspend     : Suspends the nanobox VM
  tunnel      : Displays port forward information for your apps running services
  up          : Runs 'nanobox create' and then 'nanobox deploy'
  update      : Updates the CLI to the newest available version
  upgrade     : Updates the nanobox docker images
  watch       : Watches your app for file changes
  `)

	return nil
}
