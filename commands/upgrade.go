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

// UpgradeCommand satisfies the Command interface for obtaining user info
type UpgradeCommand struct{}

// Help
func (c *UpgradeCommand) Help() {
	ui.CPrint(`
Description:
  Updates the nanobox docker images

Usage:
  pagoda upgrade
  `)
}

// Run
func (c *UpgradeCommand) Run(opts []string) {
	fmt.Printf(stylish.Bullet("Updating nanobox docker images..."))

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
	upgrade := nsync{
		model:   "imageupdate",
		path:    fmt.Sprintf("http://%v:1757/image-update", config.Nanofile.IP),
		verbose: fVerbose,
	}

	//
	upgrade.run(flags.Args())

	//
	switch upgrade.Status {

	// complete
	case "complete":
		fmt.Printf(stylish.Bullet("Upgrade complete"))

	// if the bootstrap fails the server should handle the message. If not, this can
	// be re-enabled
	case "errored":
		// fmt.Printf(stylish.Error("Bootstrap failed", "Your app failed to bootstrap"))
	}
}
