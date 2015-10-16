// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package vagrant

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/file"
)

//
var (
	err error //
)

// HaveImage returns based on wether the nanobox vagrant image is found on the
// machine or not
func HaveImage() bool {
	_, err := os.Stat(config.Root + "/nanobox-boot2docker.box")
	return err == nil
}

// runInContext runs a command in the context of a Vagrantfile (from the same dir)
func runInContext(cmd *exec.Cmd) error {

	// run the command from ~/.nanobox/apps/<config.App>. if the directory doesn't
	// exist, simply return; running the command from the directory that contains
	// the Vagratfile ensure that the command can atleast run (especially in cases
	// like 'create' where a VM hadn't been created yet, and a UUID isn't available)
	setContext(config.AppDir)

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

	// create a stdout pipe that will allow for scanning the output line-by-line;
	// if needed a stderr pipe could also be created at some point
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		config.Fatal("[util/vagrant/vagrant] cmd.StdoutPipe() failed - ", err.Error())
	}

	// create a stderr pipe that will write any error messages to the log
	stderr, err := cmd.StderrPipe()
	if err != nil {
		config.Fatal("[util/vagrant/vagrant] cmd.StderrPipe() failed - ", err.Error())
	}

	// scan the command output intercepting only 'important' lines of vagrant output'
	// and tailoring their message so as to not flood the output.
	// styled according to: http://nanodocs.gopagoda.io/engines/style-guide
	stdoutScanner := bufio.NewScanner(stdout)
	go func() {
		for stdoutScanner.Scan() {

			// log each line of output to the log
			Log.Info(stdoutScanner.Text())

			//
			switch stdoutScanner.Text() {
			case fmt.Sprintf("==> %v: VirtualBox VM is already running.", config.Nanofile.Name):
				continue
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

	// log any command errors to the log
	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			Log.Error(stderrScanner.Text())
		}
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
	setContext(config.CWDir)

	return nil
}

// add adds the nanobox vagrant image to the list of images (always overriding the
// currently installed image)
func add() error {
	return exec.Command("vagrant", "box", "add", "--force", "--name", "nanobox/boot2docker", config.Root+"/nanobox-boot2docker.box").Run()
}

// download downloads the newest nanobox vagrant image and the corresponding md5
// hash
func download() error {

	// download mv
	box, err := os.Create(config.Root + "/nanobox-boot2docker.box")
	if err != nil {
		config.Fatal("[util/vagrant/vagrant] os.Create() failed - ", err.Error())
	}
	defer box.Close()

	//
	if err := file.Progress(fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.box"), box); err != nil {
		return err
	}

	//
	// download vm md5
	md5, err := os.Create(config.Root + "/nanobox-boot2docker.md5")
	if err != nil {
		config.Fatal("[util/vagrant/vagrant] os.Create() failed - ", err.Error())
	}
	defer md5.Close()

	//
	return file.Download("https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.md5", md5)
}

// setContext changes the working directory to the designated context
func setContext(context string) {
	if err := os.Chdir(context); err != nil {
		config.Fatal("[util/vagrant/vagrant] os.Chdir() failed -", err.Error())
	}
}
