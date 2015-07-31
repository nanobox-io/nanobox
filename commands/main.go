// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// Commands/subCommands represents a map of all the available nanobox cli commands
var Commands map[string]Command

// Command represents a nanobox CLI command
type Command interface {
	Help()             // Prints the help text associated with this command
	Run(opts []string) // Houses the logic that will be run upon calling this command
}

// init builds the list of available nanobox commands and sub commands
func init() {

	// the map of all available nanobox commands
	Commands = map[string]Command{
		"bootstrap": &BootstrapCommand{},
		"build":     &BuildCommand{},
		"console":   &ConsoleCommand{},
		"create":    &CreateCommand{},
		"deploy":    &DeployCommand{},
		"destroy":   &DestroyCommand{},
		"domain":    &DomainCommand{},
		"exec":      &ExecCommand{},
		"fetch":     &FetchCommand{},
		"halt":      &HaltCommand{},
		"help":      &HelpCommand{},
		"init":      &InitCommand{},
		"log":       &LogCommand{},
		"new":       &NewCommand{},
		"publish":   &PublishCommand{},
		"reload":    &ReloadCommand{},
		"resume":    &ResumeCommand{},
		"ssh":       &SSHCommand{},
		"status":    &StatusCommand{},
		"suspend":   &SuspendCommand{},
		"tunnel":    &TunnelCommand{},
		"up":        &UpCommand{},
		"update":    &UpdateCommand{},
		"upgrade":   &UpgradeCommand{},
		"watch":     &WatchCommand{},
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

	// create a pipe that we can pipe the cmd standard output's too. The reason this
	// is done rather than just piping directly to os standard outputs and .Run()ing
	// the command (vs .Start()ing) is because the output needs to be modified
	// according to http://nanodocs.gopagoda.io/engines/style-guide
	//
	// NOTE: the reason it's done this way vs using the cmd.*Pipe's is so that all
	// the command output can be read from a single pipe, rather than having to create
	// a new pipe/scanner for each type of output
	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()

	// connect standard output
	cmd.Stdout = pw
	cmd.Stderr = pw

	//
	fmt.Printf(stylish.Bullet(fmt.Sprintf("running '%v'", strings.Trim(fmt.Sprint(cmd.Args), "[]"))))

	// scan the command output modifying it according to
	// http://nanodocs.gopagoda.io/engines/style-guide
	scanner := bufio.NewScanner(pr)
	go func() {
		for scanner.Scan() {
			fmt.Printf("   %s\n", scanner.Text())
		}
	}()

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	if err := cmd.Start(); err != nil {
		ui.LogFatal("[commands.runVagrantCommand] cmd.Start() failed", err)
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		ui.LogFatal("[commands.runVagrantCommand] cmd.Wait() failed", err)
	}

	// switch back to project dir
	if err := os.Chdir(config.CWDir); err != nil {
		ui.LogFatal("[commands.runVagrantCommand] os.Chdir() failed", err)
	}
}
