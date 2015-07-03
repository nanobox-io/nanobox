// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/kardianos/osext"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

// UpdateCommand satisfies the Command interface for obtaining user info
type UpdateCommand struct{}

// Help prints detailed help text for the user command
func (c *UpdateCommand) Help() {
	ui.CPrint(`
Description:
  Update the CLI to the most recently released version.

Usage:
  pagoda update

  ex. pagoda update
  `)
}

// Run gets the current user and prints out select information to the terminal
func (c *UpdateCommand) Run(opts []string) {

	fmt.Println("Updating...")

	//
	program := os.Args[0]

	//
	config.Console.Info("[commands.update] Program: %v", program)

	//
	path, err := osext.Executable()
	if err != nil {
		ui.LogFatal("[commands.update] osext.ExecutableFolder() failed", err)
	}

	//
	config.Console.Info("[commands.update] Path: %v", path)

	// download a new CLI from s3 that matches their os and arch
	download := fmt.Sprintf("https://s3-us-west-2.amazonaws.com/tools.nanobox.io/cli/%v/%v/nanobox", runtime.GOOS, runtime.GOARCH)

	//
	config.Console.Info("[commands.update] Downloading new CLI from %v", download)

	// create a new request
	req, err := http.NewRequest("GET", download, nil)
	if err != nil {
		ui.LogFatal("[commands.update] http.NewRequest() failed", err)
	}

	// download the new CLI
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ui.LogFatal("[commands.update] http.DefaultClient.Do() failed", err)
	}

	var buf bytes.Buffer
	p := make([]byte, 2048)
	total := 0

	// replace the existing CLI with the new CLI
	for {

		// stream read the download
		n, err := res.Body.Read(p)

		// write to our buffer
		buf.Write(p[:n])

		// update the total bytes read
		total += n

		// show how download progress
		fmt.Printf("\rDownloading... %d/%d", total, res.ContentLength)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				ui.LogFatal("[commands.update] res.Body.Read() failed", err)
			}
		}

		defer res.Body.Close()
	}
	ui.CPrint("\rDownloading... [green]success[reset]                           ")

	//
	fmt.Println("Writing to", path)
	ioutil.WriteFile(path, buf.Bytes(), 0755)

	//
	fmt.Println("Update successful!")

	//
	// config.Console.Debug("[commands.update] command: %v, args: %+v", os.Args[0], os.Args[1:])

	// // attempt to run the command that was being run to begin with (unless its update)
	// if os.Args[1] != "update" {

	//  // run the command that was being run to begin with
	//  out, err := exec.Command(os.Args[0], os.Args[1:]...).Output()
	//  if err != nil {
	//    ui.LogFatal("[commands.update] exec.Command()", err)
	//  }

	//  // show the output of the command that is run
	//  fmt.Printf("%v\n", out)
	// }
}
