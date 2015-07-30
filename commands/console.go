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
	"os"

	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// ConsoleCommand satisfies the Command interface
type ConsoleCommand struct{}

// Help
func (c *ConsoleCommand) Help() {
	ui.CPrint(`
Description:
  Drops you into bash inside of the nanobox vm docker

Usage:
  nanobox console
  `)
}

// Run
func (c *ConsoleCommand) Run(opts []string) {
	fmt.Printf(stylish.Bullet("Opening a nanobox console..."))

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	//
	var fTunnel string
	flags.StringVar(&fTunnel, "t", "", "")
	flags.StringVar(&fTunnel, "tunnel", "", "")

	//
	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.destroy] flags.Parse() failed", err)
	}

	//
	if len(flags.Args()) > 0 {
		fmt.Println("Attempting to run 'nanobox console' with a command. Use 'nanobox exec'")
		os.Exit(0)
	}

	// add a check here to regex the fTunnel to make sure the format makes sense

	//
	v := url.Values{}
	v.Add("forward", fTunnel)

	console := utils.Docker{Params: v.Encode()}
	console.Run()
}
