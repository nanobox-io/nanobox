// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"github.com/spf13/cobra"

	// "github.com/pagodabox/nanobox-cli/config"
	// "github.com/pagodabox/nanobox-cli/util"
	// "github.com/pagodabox/nanobox-golang-stylish"
)

//
var boxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "",
	Long:  ``,

	Run: boxUpdate,
}

// boxUpdate
func boxUpdate(ccmd *cobra.Command, args []string) {
	// Install()
	// 	release := latestVersion()
	// 	if currentVersion() >= release.version() {
	// 		fmt.Println("I already have the latest")
	// 		return
	// 	}
	// 	// asset := release.Assets[0]
	// 	// put file downloader here downloading from asset.DownloadURL
	// 	setVersion(release.version())
	// 	// vagrant box add ~/.nanobox/boot2docker.box
}
