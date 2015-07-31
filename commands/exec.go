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
	"strings"

	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// ExecCommand satisfies the Command interface
type ExecCommand struct{}

// Help
func (c *ExecCommand) Help() {
	ui.CPrint(`
Description:
  Runs a command from inside your app on the nanobox VM

Usage:
  nanobox exec [-t 80:8080] <command>

Options:
	-t, --tunnel
		Creates port forwards for all comma delimeted port:port combos
	`)
}

// Run
func (c *ExecCommand) Run(opts []string) {
	fmt.Printf(stylish.Bullet("Running command..."))

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
	cmd := flags.Args()[0:]

	//
	if len(cmd) <= 0 {
		cmd = append(cmd, ui.Prompt("Please specify a command you wish to exec: "))
	}

	// add a check here to regex the fTunnel to make sure the format makes sense

	//
	v := url.Values{}
	v.Add("forward", fTunnel)
	v.Add("cmd", strings.Join(cmd, " "))

	exec := utils.Docker{Params: v.Encode()}
	exec.Run()
}
