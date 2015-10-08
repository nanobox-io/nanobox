// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-cli/util/file"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the CLI to the newest available version",
	Long:  ``,

	Run: update,
}

// update
func update(ccmd *cobra.Command, args []string) {

	// if the local md5 matches the remote md5 there is no need to update
	match, err := util.MD5sMatch(config.Root+"/nanobox.md5", "https://s3.amazonaws.com/tools.nanobox.io/cli/nanobox.md5")
	if err != nil {
		config.Fatal("[commands/update] util.MD5sMatch() failed", err.Error())
	}

	if match {
		fmt.Printf("Nanobox is up to date (running v%s)", config.VERSION)
		return
	}

	fmt.Printf(stylish.Bullet("Updating nanobox"))

	//
	// get the path of the current executing CLI
	path, err := osext.Executable()
	if err != nil {
		config.Fatal("[commands/update] osext.ExecutableFolder() failed", err.Error())
	}

	// download the CLI
	cli, err := os.Create(path)
	if err != nil {
		config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer cli.Close()

	//
	file.Progress(fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/nanobox", runtime.GOOS, runtime.GOARCH), cli)

	//
	// download the CLI md5
	md5, err := os.Create(config.Root + "/nanobox.md5")
	if err != nil {
		config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer md5.Close()

	//
	file.Download("https://s3.amazonaws.com/tools.nanobox.io/cli/nanobox.md5", md5)

	// if the new CLI fails to execute, just print a generic message and return
	out, err := exec.Command(path, "-v").Output()
	if err != nil {
		fmt.Printf(stylish.SubBullet("[√] Update successful"))
		return
	}

	//
	fmt.Printf(stylish.SubBullet("[√] Now running v%s", string(out)))
}
