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
	"strings"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// GetVMUUID tries to return the VMs uuid found in it's corresponding .vangrant
// folder. If a uuid is not found than the VM has not yet been created.
func GetVMUUID() ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("%v/.vagrant/machines/%v/%v/index_uuid", config.AppDir, config.App, config.Nanofile.Provider))
}

//
func GetVMStatus() error {

	//
	cmd := exec.Command("vagrant", "global-status")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	//
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanRunes)
	go func() {
		for scanner.Scan() {

			fmt.Println(scanner.Text())

			// print line
			// switch scanner.Text() {
			// case "\n", "\r":
			// 	fmt.Printf("%s   ", scanner.Text())
			// default:
			// 	fmt.Print(scanner.Text())
			// }
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

	return nil
}

// RunVagrantCommand provides a wrapper around a standard cmd.Run() in which
// all standard in/outputs are connected to the command, and the directory is
// changed to the corresponding app directory. This allows nanobox to run Vagrant
// commands w/o contaminating a users codebase.
func RunVagrantCommand(cmd *exec.Cmd) error {

	// GetVMStatus()
	// os.Exit(0)

	fmt.Printf(stylish.Bullet("running '%v'", strings.Trim(fmt.Sprint(cmd.Args), "[]")))

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

	// scan the command output modifying it according to
	// http://nanodocs.gopagoda.io/engines/style-guide
	scanner := bufio.NewScanner(pr)
	scanner.Split(bufio.ScanRunes)
	go func() {
		for scanner.Scan() {

			// print line
			switch scanner.Text() {
			case "\n", "\r":
				fmt.Printf("%s   ", scanner.Text())
			default:
				fmt.Print(scanner.Text())
			}
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
	if err := os.Chdir(config.CWDir); err != nil {
		return err
	}

	return nil
}
