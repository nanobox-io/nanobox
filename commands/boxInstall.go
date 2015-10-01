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
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var boxInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,

	Run: boxInstall,
}

// boxInstall
func boxInstall(ccmd *cobra.Command, args []string) {

	//
	boxfile := filepath.Clean(config.Root + "/nanobox-boot2docker.box")

	//
	if _, err := os.Stat(boxfile); err != nil {
		fmt.Printf(stylish.Bullet("Installing nanobox image..."))

		//
		util.VMDownload()

		// always replace the existing box with the new box
		if err := exec.Command("vagrant", "box", "add", "--force", "--name", "nanobox/boot2docker", boxfile).Run(); err != nil {
			config.Fatal("[commands/boxInstall]", err.Error())
		}
	}
}
