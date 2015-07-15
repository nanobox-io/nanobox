// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// Commands represents a map of all the available commands that the Pagoda Box
// CLI can run
var Commands map[string]Command

// Command represents a Pagoda Box CLI command. Every command must have a Help()
// and Run() function
type Command interface {
	Help()             // Prints the help text associated with this command
	Run(opts []string) // Houses the logic that will be run upon calling this command
}

// init builds the list of available Pagoda Box CLI commands
func init() {

	// the map of all available commands the Pagoda Box CLI can run
	Commands = map[string]Command{
		"build":   &BuildCommand{},
		"create":  &CreateCommand{},
		"domain":  &DomainCommand{},
		"deploy":  &DeployCommand{},
		"destroy": &DestroyCommand{},
		"fetch":   &FetchCommand{},
		"halt":    &HaltCommand{},
		"help":    &HelpCommand{},
		"init":    &InitCommand{},
		"log":     &LogCommand{},
		"new":     &NewCommand{},
		"publish": &PublishCommand{},
		"reload":  &ReloadCommand{},
		"resume":  &ResumeCommand{},
		"ssh":     &SSHCommand{},
		"status":  &StatusCommand{},
		"suspend": &SuspendCommand{},
		"up":      &UpCommand{},
		"update":  &UpdateCommand{},
	}
}

// runVagrantCommand provides a wrapper around a standard cmd.Run() in which
// all standard in/outputs are connected to the command, and the directory is
// changed to the corresponding app directory. This allows nanobox to run Vagrant
// commands w/o contaminating a users codebase
func runVagrantCommand(cmd *exec.Cmd) {

	// run an init to ensure there is a Vagrantfile
	init := InitCommand{}
	init.Run(nil)

	// run the command from ~/.nanobox/apps/<this app>
	if err := os.Chdir(config.AppDir); err != nil {
		ui.LogFatal("[commands.runVagrantCommand] os.Chdir() failed", err)
	}

	// connect standard in/outputs
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//
	fmt.Printf(stylish.Bullet(fmt.Sprintf("running '%v'", strings.Trim(fmt.Sprint(cmd.Args), "[]"))))

	// run command; if it fails Vagrant will output something and we'll just exit
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	// switch back to project dir
	if err := os.Chdir(config.CWDir); err != nil {
		ui.LogFatal("[commands.runVagrantCommand] os.Chdir() failed", err)
	}
}
