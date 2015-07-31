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

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// BootstrapCommand satisfies the Command interface for deploying to nanobox
type BootstrapCommand struct{}

// Help
func (c *BootstrapCommand) Help() {
	ui.CPrint(`
Description:
  Runs an engine's bootstrap script - downloads code & launches VM

Usage:
  nanobox bootstrap [-v]

Options:
  -v, --verbose
    Increase the level of log output from 'info' to 'debug'
  `)
}

// Run issues a deploy to the running nanobox VM
func (c *BootstrapCommand) Run(opts []string) {
	fmt.Printf(stylish.Bullet("Bootstrapping code..."))

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	//
	var fVerbose bool
	flags.BoolVar(&fVerbose, "v", false, "")
	flags.BoolVar(&fVerbose, "verbose", false, "")

	//
	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.destroy] flags.Parse() failed", err)
	}

	//
	bootstrap := nsync{
		kind:    "bootstrap",
		path:    fmt.Sprintf("http://%v:1757/bootstrap", config.Nanofile.IP),
		verbose: fVerbose,
	}

	//
	bootstrap.run(flags.Args())
}
