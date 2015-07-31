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
	"runtime"

	"github.com/kardianos/osext"

	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// UpdateCommand satisfies the Command interface for obtaining user info
type UpdateCommand struct{}

// Help
func (c *UpdateCommand) Help() {
	ui.CPrint(`
Description:
  Updates the CLI to the newest available version

Usage:
  pagoda update
  `)
}

// Run
func (c *UpdateCommand) Run(opts []string) {

	fmt.Printf(stylish.Bullet("Updating nanobox CLI..."))

	//
	path, err := osext.Executable()
	if err != nil {
		ui.LogFatal("[commands.update] osext.ExecutableFolder() failed", err)
	}

	//
	fmt.Printf(stylish.SubBullet(fmt.Sprintf("Nanobox CLI found running at %v", path)))

	// download a new CLI from s3 that matches their os and arch
	download := fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/nanobox", runtime.GOOS, runtime.GOARCH)

	// create a new request
	fmt.Printf(stylish.SubBullet(fmt.Sprintf("Downloading latest CLI from %v", download)))
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
	ui.CPrint("\rDownloading... [green]success[reset]")

	//
	fmt.Printf(stylish.SubBullet(fmt.Sprintf("Replacing CLI at %v", path)))
	ioutil.WriteFile(path, buf.Bytes(), 0755)

	//
	fmt.Println(stylish.Success())

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
