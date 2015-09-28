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
	// "io"
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

	// if no uuid is found the vm has not yet been created
	if uuid == "" {
		return "not created"
	}

	// run the command with the uuid; this allows the command to be run from any
	// context (not dependent on a Vagrantfile)
	cmd := exec.Command("vagrant", "status", uuid)

	//
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		config.Fatal("[util/vagrant] cmd.StdoutPipe() failed", err.Error())
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
				submatch := regexp.MustCompile(`\s(running|saved|poweroff|not created|aborted)\s`).FindStringSubmatch(scanner.Text())

				//
				if len(submatch) > 1 {
					status = submatch[1]
				} else {
					config.Fatal("[util/vagrant] Unknown status: ", scanner.Text())
				}
			}
		}
	}()

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	if err := cmd.Start(); err != nil {
		config.Fatal("[util/vagrant] cmd.Start() failed", err.Error())
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		config.Fatal("[util/vagrant] cmd.Wait() failed", err.Error())
	}

	// set the status in the .vmfile (currently not used)
	// config.VMfile.StatusIs(status)

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

	// create a stdout pipe that will allow for scanning the output line-by-line;
	// if needed a stderr pipe could also be created at some point
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		config.Fatal("[util/vagrant] cmd.StdoutPipe() failed", err.Error())
	}

	// start a goroutine that will act as an 'outputer' allowing us to add 'dots'
	// to the end of each line (as these lines are a reduced version of the actual
	// output there will be some delay between output)
	output := make(chan string)
	go func() {

		tick := time.Second

		// block until any one message outputs
		msg, ok := <-output

		// print initial message to 'get the ball rolling' on our 'outputer'
		fmt.Printf("   - %s", msg)

		// begin a loop to read off the channel until it's closed
		for {
			select {

			// print any messages and reset ticker
			case msg, ok = <-output:

				// once the channel closes print the final newline and close the goroutine
				if !ok {
					fmt.Println("")
					return
				}

				fmt.Printf("\n   - %s", msg)

				tick = time.Second

			// after every tick print a '.' until we get another message one the channel
			// (at which point ticker is reset and it starts all over again)
			case <-time.After(tick):
				fmt.Print(".")

				// increase the wait time by half of the total previous time
				tick += tick / 2
			}
		}
	}()

	// scan the command output intercepting only 'important' lines of vagrant output'
	// and tailoring their message so as to not flood the output.
	// styled according to: http://nanodocs.gopagoda.io/engines/style-guide
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {

			//
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
			// case fmt.Sprintf("==> %v: Destroying VM and associated drives...", config.Nanofile.Name):
			// 	output <- "Destroying virtual machine"
			case fmt.Sprintf("==> %v: Forcing shutdown of VM...", config.Nanofile.Name):
				output <- "Shutting down virtual machine"
			case fmt.Sprintf("==> %v: Saving VM state and suspending execution...", config.Nanofile.Name):
				output <- "Saving virtual machine"
				// case fmt.Sprintf("==> %v: Resuming suspended VM...", config.Nanofile.Name):
				// 	output <- "Resuming virtual machine"
			}
		}

		// close the output channel once all lines of command output have been read
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
