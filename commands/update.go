// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package commands

import (
	"fmt"
	"github.com/kardianos/osext"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	fileutil "github.com/nanobox-io/nanobox/util/file"
	printutil "github.com/nanobox-io/nanobox/util/print"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"runtime"
	"time"
)

//
var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Updates the CLI to the newest available version",
		Long:  ``,

		Run: update,
	}

	devv bool // update the dev version of the CLI
)

//
func init() {
	updateCmd.Flags().BoolVarP(&devv, "dev", "d", false, "Update dev CLI")
}

// update
func update(ccmd *cobra.Command, args []string) {

	cli, match, err := getUpdateStuff()
	if err != nil {
		Config.Fatal("[commands/update] getUpdateStuff() failed", err.Error())
	}

	// if the md5s don't match or it's been forced, update
	switch {
	case config.Force, !match:
		runUpdate(cli, devv)
	default:
		fmt.Printf(stylish.SubBullet("[√] Nanobox is up-to-date"))
	}
}

// Update
func Update() {

	cli, match, err := getUpdateStuff()
	if err != nil {
		Config.Fatal("[commands/update] getUpdateStuff() failed", err.Error())
	}

	// stat the update file to get ModTime(); an error here means the file doesn't
	// exist. This is highly unlikely as the file is created if it doesn't exist
	// each time the CLI is run.
	fi, _ := os.Stat(config.UpdateFile)

	// if the md5s don't match and it's 'time' for an update (14 days), OR a force
	// update is issued, update
	if !match && time.Since(fi.ModTime()).Hours() >= 336.0 {

		//
		switch printutil.Prompt("Nanobox is out of date, would you like to update it now (y/N)? ") {

		// don't update by default
		default:
			fmt.Println("You can manually update at any time with 'nanobox update'.")
			return

		// if yes continue to update
		case "Yes", "yes", "Y", "y":
			runUpdate(cli, false)
		}
	}
}

// runUpdate
func runUpdate(path string, devv bool) {

	fmt.Printf(stylish.Bullet("Updating nanobox..."))

	// create a new CLI at the location of the old one
	cli, err := os.Create(path)
	if err != nil {
		Config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer cli.Close()

	download := fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/nanobox", runtime.GOOS, runtime.GOARCH)
	if devv {
		download += "-dev"
	}

	// download the new cli
	fileutil.Progress(download, cli)

	// ensure the newly downloaded cli matches the remote
	match, err := Util.MD5sMatch(path, fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/nanobox.md5", runtime.GOOS, runtime.GOARCH))
	if err != nil {
		Config.Fatal("[commands/update] util.MD5sMatch() failed", err.Error())
	}

	// if they don't match it's the wrong CLI
	if !match {
		fmt.Println("MD5 checksum failed! Your nanobox-desktop (CLI) is not ours!")
		return
	}

	// if the new CLI fails to execute, just print a generic message and return
	out, err := exec.Command(path, "-v").Output()
	if err != nil {
		fmt.Printf(stylish.SubBullet("[√] Update successful"))
		return
	}

	fmt.Printf(stylish.SubBullet("[√] Now running %s", string(out)))

	// update the modification time of the .update file
	if err := os.Chtimes(config.UpdateFile, time.Now(), time.Now()); err != nil {
		Config.Fatal("[commands.update] os.Chtimes() failed", err.Error())
	}
}

// getUpdateStuff
func getUpdateStuff() (cli string, match bool, err error) {

	// get the path of the current executing CLI
	if cli, err = osext.Executable(); err != nil {
		return
	}

	// check the current cli md5 against the remote md5
	if match, err = Util.MD5sMatch(cli, "https://s3.amazonaws.com/tools.nanobox.io/cli/nanobox.md5"); err != nil {
		return
	}

	return
}
