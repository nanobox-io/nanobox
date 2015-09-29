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
	"runtime"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the CLI to the newest available version",
	Long: `
Description:
  Updates the CLI to the newest available version`,

	Run: nanoUpdate,
}

// nanoUpdate
func nanoUpdate(ccmd *cobra.Command, args []string) {

	fmt.Printf(stylish.Bullet("Updating nanobox CLI"))

	//
	md5, err := os.Create(config.Root + "/nanobox.md5")
	if err != nil {
		config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer md5.Close()

	// download the cli md5
	util.Download("https://s3.amazonaws.com/tools.nanobox.io/cli/nanobox.md5", md5)

	//
	// get the path of the current executing CLI
	path, err := osext.Executable()
	if err != nil {
		config.Fatal("[commands/update] osext.ExecutableFolder() failed", err.Error())
	}

	//
	cli, err := os.Create(path)
	if err != nil {
		config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer cli.Close()

	// download the CLI with a progres bar
	util.Progress(fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/nanobox", runtime.GOOS, runtime.GOARCH), cli)

	//
	fmt.Printf(stylish.SubBullet("[√] Now running %s", config.VERSION))
}
