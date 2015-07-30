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
	"net/url"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// DeployCommand satisfies the Command interface for deploying to nanobox
type DeployCommand struct{}

// Help
func (c *DeployCommand) Help() {
	ui.CPrint(`
Description:
  Issues a deploy to your nanobox

Usage:
  nanobox deploy
  nanobox deploy -v
  nanobox deploy -r

Options:
  -v, --verbose
    Increase the level of log output from 'info' to 'debug'

  -r, --reset
    Clears cached libraries the project might use
  `)
}

// Run issues a deploy to the running nanobox VM
func (c *DeployCommand) Run(opts []string) {
	fmt.Printf(stylish.Bullet("Deploying codebase..."))

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	//
	var fReset bool
	flags.BoolVar(&fReset, "r", false, "")
	flags.BoolVar(&fReset, "reset", false, "")

	//
	var fVerbose bool
	flags.BoolVar(&fVerbose, "v", false, "")
	flags.BoolVar(&fVerbose, "verbose", false, "")

	//
	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.destroy] flags.Parse() failed", err)
	}

	v := url.Values{}

	if fReset {
		v.Add("reset", "true")
	}

	//
	deploy := nsync{
		kind:    "deploy",
		path:    fmt.Sprintf("http://%v:1757/deploys?%v", config.Nanofile.IP, v.Encode()),
		verbose: fVerbose,
	}

	//
	deploy.run(flags.Args())
}
