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
	// "github.com/pagodabox/nanobox-golang-stylish"
)

// ExecCommand satisfies the Command interface
type ExecCommand struct{}

// Help
func (c *ExecCommand) Help() {
	ui.CPrint(`
Description:
  Run's a command on an application's service. 'compound' commands MUST be
  formated inside single quotes ('')

    'ls -la'

Usage:
  nanobox run <COMMAND>
  nanobox run '<COMPOUND COMMAND>'

	ex. pagoda run ls
  ex. pagoda run 'ls -la'
	`)
}

// Run
func (c *ExecCommand) Run(opts []string) {

	// If there's no command, prompt for one
	if len(opts) <= 0 {
		opts[0] = ui.Prompt("Please specify a command you wish to run (see help for format): ")
	}

	// if there are too many options they probably forgot that a compound command
	// needs to be a string
	if len(opts) > 1 {
		fmt.Println("Too many arguments... Please ensure you place compound commands in single quotes ('')")
	}

	console := ConsoleCommand{}
	console.Run(opts)
}
