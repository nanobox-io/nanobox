// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

// NewCommand satisfies the Command interface
type NewCommand struct{}

// Help prints detailed help text for the app list command
func (c *NewCommand) Help() {
	ui.CPrint(`
Description:
  Generate a new nanobox project in the current working directory

  type:
    the type of project (engine/plugin/service)

  name:
    the name of the project

Usage:
  nanobox new <type> <name>

  ex. nanobox new engine nodejs

Options:
  -m, --minimal
    Create a minimal version
  `)
}

// Run destroys the specified virtual machine
func (c *NewCommand) Run(opts []string) {

	fmt.Println("OTPS!", opts)

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fStyle bool
	flags.BoolVar(&fStyle, "m", false, "")
	flags.BoolVar(&fStyle, "minimal", false, "")

	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.new] flags.Parse() failed", err)
	}

	if len(opts) < 2 {
		fmt.Println("NOT ENOUGH ARGS!")
		os.Exit(1)
	}

	typ := opts[0]
	name := fmt.Sprintf("nanobox-%s", opts[1])
	version := "0.0.1"

	// create a new project by the name, unless it already exists
	if di, _ := os.Stat(name); di == nil {

		//
		config.Console.Info("[install] Creating '%v' directory", name)

		if err := os.Mkdir(name, 0755); err != nil {
			fmt.Println("BONK!")
		}

		entry := fmt.Sprintf(`
name: %-18s     # the name of your project
type: %-18s     # the type of project (engine/plugin/service)
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
includes:
  -
`, name, typ, version)

		config.Console.Info("[install] Creating Packagefile...")

		path := name + "/Packagefile"

		if _, err := os.Create(path); err != nil {
			fmt.Println("BONK!")
		}

		// write the Packagefile
		if err := ioutil.WriteFile(path, []byte(entry), 0644); err != nil {
			fmt.Println("BONK!", err)
		}

	} else {
		fmt.Printf("A project by the name '%s' already exists at this location...\n", name)
	}

}
