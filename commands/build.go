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

// BuildCommand satisfies the Command interface for deploying to nanobox
type BuildCommand struct{}

// Help
func (c *BuildCommand) Help() {
	ui.CPrint(`
Description:
  Rebuilds/compiles your project

Usage:
  nanobox build [-v]

Options:
  -v, --verbose
    Increase the level of log output from 'info' to 'debug'
  `)
}

// Run issues a deploy to the running nanobox VM
func (c *BuildCommand) Run(opts []string) {
	fmt.Printf(stylish.Bullet("Building codebase..."))

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
	build := nsync{
		kind:    "build",
		path:    fmt.Sprintf("http://%v:1757/builds", config.Nanofile.IP),
		verbose: fVerbose,
	}

	//
	build.run(flags.Args())

	//
	switch build.Status {

	// complete
	case "complete":
		fmt.Printf(stylish.Bullet(fmt.Sprintf("Build complete... Navigate to %v.nano.dev to view your app.", config.App)))

		// if the build fails the server should handle the message. If not, this can
		// be re-enabled
	case "errored":
		// fmt.Printf(stylish.Error("Build failed", "Your build failed to well... build"))
	}
}
