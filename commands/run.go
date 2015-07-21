// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	// "bufio"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	// "io"
	"os"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// RunCommand satisfies the Command interface
type RunCommand struct{}

// Help prints detailed help text for the app list command
func (c *RunCommand) Help() {
	ui.CPrint(`
Description:
  Drops you into bash inside your nanobox vm

Usage:
  nanobox enter
  `)
}

// Run resumes the specified virtual machine
func (c *RunCommand) Run(opts []string) {

	// create an SSH client
	client, err := ssh.Dial("tcp", config.Boxfile.IP+":22", &ssh.ClientConfig{User: "docker", Auth: []ssh.AuthMethod{ssh.Password("tcuser")}})
	if err != nil {
		ui.LogFatal("[commands.service_ssh] ssh.Dial() failed", err)
	}
	defer client.Close()

	// create an SSH session for the client
	session, err := client.NewSession()
	if err != nil {
		ui.LogFatal("[commands.service_ssh] client.NewSession() failed", err)
	}
	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 24, 80, modes); err != nil {
		ui.LogFatal("[commands.service_ssh] session.RequestPty() failed", err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	args := ""

	// iterate over the args (minus program and command)
	for _, a := range os.Args[2:] {
		args += fmt.Sprintf("%s ", a)
	}

	cmd := fmt.Sprintf(`
docker \
  run \
    -it \
    --rm \
    -v /mnt/sda/var/nanobox/deploy/:/data/ \
    -v /vagrant/code/nanobox-ruby-sample/:/code/ \
    -w /code \
    -e PATH=/data/sbin:/data/bin:/opt/gonano/sbin:/opt/gonano/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    nanobox/build \
    %v`, args)

	fmt.Printf(stylish.Bullet(fmt.Sprintf("Running command: %v", args)))

	// run a command
	if err := session.Run(cmd); err != nil {
		ui.LogFatal("Failed to run command", err)
	}
}
