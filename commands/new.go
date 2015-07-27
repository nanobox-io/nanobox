// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	// "github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// NewCommand satisfies the Command interface
type NewCommand struct{}

// Help
func (c *NewCommand) Help() {
	ui.CPrint(`
Description:
  Generate a new nanobox project in the current working directory

  name:
    the name of the project

Usage:
  nanobox new <name>

  ex. nanobox new ruby
  `)
}

// Run
func (c *NewCommand) Run(opts []string) {

	if len(opts) < 1 {
		fmt.Println(stylish.Error("app name required", "Please provide a name when generating a new engine, (run 'nanobox new -h' for details)"))
		os.Exit(1)
	}

	name := fmt.Sprintf("nanobox-%s", opts[0])
	version := "0.0.1"

	fmt.Println(stylish.Header("initializing new app: " + name))

	// create a new project by the name, unless it already exists
	if di, _ := os.Stat(name); di == nil {

		stylish.Bullet("Creating engine at")

		if err := os.Mkdir(name, 0755); err != nil {
			ui.LogFatal("[commands.new] os.Mkdir() failed", err)
		}

		entry := fmt.Sprintf(`
name: %-18s     # the name of your project
version: %-18s  # the current version of the project
summary:                     # a short summary of the project
description:                 # a detailed description of the project
stability: alpha             # the current stability of the project (alpha/beta/stable)
license: MIT                 # the license to be applied to this software
readme: README.md            # the path/to/your/readme.file

# a list of authors/contributors
authors:
  -

# the list of files your project requires. These are bundled
# into a 'release.tar.gz' and uploaded to nanobox
project_files:
  -
`, name, version)

		path := name + "/Enginefile"

		if _, err := os.Create(path); err != nil {
			ui.LogFatal("[commands.new] os.Create() failed", err)
		}

		// write the Enginefile
		if err := ioutil.WriteFile(path, []byte(entry), 0644); err != nil {
			ui.LogFatal("[commands.new] ioutil.WriteFile() failed", err)
		}

	} else {
		fmt.Printf("A project by the name '%s' already exists at this location...\n", name)
	}

	stylish.Bullet("Default Enginefile created at")
}
