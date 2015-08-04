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
	"strconv"

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
  Issues a deploy to the nanobox VM

Usage:
  nanobox deploy [-v] [-r] [-s]

Options:
  -v, --verbose
    Increases the level of log output from 'info' to 'debug'

  -r, --reset
    Clears cached libraries the project might use

	-s, --sandbox
		Creates your app environment w/o webs or workers
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
	var fSandbox bool
	flags.BoolVar(&fSandbox, "s", false, "")
	flags.BoolVar(&fSandbox, "sandbox", false, "")

	//
	var fVerbose bool
	flags.BoolVar(&fVerbose, "v", false, "")
	flags.BoolVar(&fVerbose, "verbose", false, "")

	//
	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.destroy] flags.Parse() failed", err)
	}

	v := url.Values{}

	v.Add("reset", strconv.FormatBool(fReset))
	v.Add("sandbox", strconv.FormatBool(fSandbox))

	//
	deploy := nsync{
		kind:    "deploy",
		path:    fmt.Sprintf("http://%v:1757/deploys?%v", config.Nanofile.IP, v.Encode()),
		verbose: fVerbose,
	}

	//
	deploy.run(flags.Args())

	//
	switch deploy.Status {

	// complete
	case "complete":

		//
		if fSandbox {
			fmt.Printf(stylish.Bullet("Sandbox deploy complete..."))
			break
		}

		fmt.Printf(stylish.Bullet(fmt.Sprintf("Deploy complete... Navigate to %v.nano.dev to view your app.", config.App)))

		// if the deploy fails the server should handle the message. If not, this can
		// be re-enabled
	case "errored":
		// fmt.Printf(stylish.Error("Deploy failed", "Your deploy failed to well... deploy"))
	}
}
