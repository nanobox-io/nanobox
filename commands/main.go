// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"os"
	"os/exec"

	"github.com/pagodabox/nanobox-cli/config"
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
func runVagrantCommand(cmd *exec.Cmd) error {

	// run the command from ~/.nanobox/apps/<this app>
	if err := os.Chdir(config.AppDir); err != nil {
		return err
	}

	// connect standard in/outputs
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run command
	if err := cmd.Run(); err != nil {
		return err
	}

	// switch back to project dir
	if err := os.Chdir(config.CWDir); err != nil {
		return err
	}

	return nil
}
