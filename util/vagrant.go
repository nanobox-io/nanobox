// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/pagodabox/nanobox-cli/config"
	// "github.com/pagodabox/nanobox-golang-stylish"
)

// GetVMUUID tries to return the VMs uuid found in it's corresponding .vangrant
// folder. If a uuid is not found than the VM has not yet been created. Don't
// really care about the error here since the value will be "" if there is no
// file to read
func GetVMUUID() string {
	b, _ := ioutil.ReadFile(fmt.Sprintf("%v/.vagrant/machines/%v/%v/index_uuid", config.AppDir, config.Nanofile.Name, config.Nanofile.Provider))
	return string(b)
}

// GetVMStatus returns the current status of the VM; this command needs to be run
// in a way independant of a Vagrantfile to ensure that the status will always
// be available
func GetVMStatus() string {

	var status string

	//
	uuid := GetVMUUID()

	if uuid == "" {
		return uuid
	}

	// run the command with the uuid; this allows the command to be run from any
	// context (not dependent on a Vagrantfile)
	cmd := exec.Command("vagrant", "status", uuid)

	//
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	//
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {

			// if the scan line starts with the name of the VM thats the line thats going
			// to have the status in it
			if strings.HasPrefix(scanner.Text(), config.Nanofile.Name) {

				// pull the status out of the line; it's easiest to just pull the exact
				// states out of the string since there aren't that many
				status = regexp.MustCompile(`\s(running|saved|poweroff|not created)\s`).FindStringSubmatch(scanner.Text())[1]
			}
		}
	}()

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		panic(err)
	}

	return status
}

// RunVagrantCommand provides a wrapper around a standard cmd.Run() in which
// all standard in/outputs are connected to the command, and the directory is
// changed to the corresponding app directory. This allows nanobox to run Vagrant
// commands w/o contaminating a users codebase.
func RunVagrantCommand(cmd *exec.Cmd) error {

	// run the command from ~/.nanobox/apps/<config.App>. if the directory doesn't
	// exist, simply return; running the command from the directory that contains
	// the Vagratfile ensure that the command can atleast run (especially in cases
	// like 'create' where a VM hadn't been created yet, and a UUID isn't available)
	if err := os.Chdir(config.AppDir); err != nil {
		return err
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
	output := make(chan string)

	//
	go func() {

		tick := time.Second

		for {
			select {

			//
			case msg, ok := <-output:

				//
				if !ok {
					fmt.Println("")
					return
				}

				//
				fmt.Printf("\n   - %s", msg)
				tick = time.Second

			//
			case <-time.After(tick):
				fmt.Print(".")

				// increase the wait time by half of the total previous time
				tick += tick / 2
			}
		}
	}()

	// scan the command output modifying it according to
	// http://nanodocs.gopagoda.io/engines/style-guide
	scanner := bufio.NewScanner(pr)
	go func() {
		for scanner.Scan() {

			// intercept and modify only 'important' lines of vagrant output' so as to
			// not flood the output
			switch scanner.Text() {
			case fmt.Sprintf("==> %v: Importing base box 'nanobox/boot2docker'...", config.Nanofile.Name):
				output <- "Importing nanobox base image"
			case fmt.Sprintf("==> %v: Booting VM...", config.Nanofile.Name):
				output <- "Booting virtual machine"
			case fmt.Sprintf("==> %v: Configuring and enabling network interfaces...", config.Nanofile.Name):
				output <- "Configuring virtual network"
			case fmt.Sprintf("==> %v: Mounting shared folders...", config.Nanofile.Name):
				output <- fmt.Sprintf("Mounting source code (%s)", config.CWDir)
			case fmt.Sprintf("==> %v: Waiting for nanobox server...", config.Nanofile.Name):
				output <- "Starting nanobox server"
			case fmt.Sprintf("==> %v: Attempting graceful shutdown of VM...", config.Nanofile.Name):
				output <- "Shutting down virtual machine"
			case fmt.Sprintf("==> %v: Destroying VM and associated drives...", config.Nanofile.Name):
				output <- "Destroying virtual machine"
			case fmt.Sprintf("==> %v: Forcing shutdown of VM...", config.Nanofile.Name):
				output <- "Shutting down virtual machine"
			case fmt.Sprintf("==> %v: Saving VM state and suspending execution...", config.Nanofile.Name):
				output <- "Suspending virtual machine"
			case fmt.Sprintf("==> %v: Resuming suspended VM...", config.Nanofile.Name):
				output <- "Resuming virtual machine"
			}
		}

		// close the channel
		close(output)
	}()

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	if err := cmd.Start(); err != nil {
		return err
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		return err
	}

	// switch back to project dir
	if err := os.Chdir(config.CWDir); err != nil {
		return err
	}

	return nil
}
