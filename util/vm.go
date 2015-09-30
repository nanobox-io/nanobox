// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

//
import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox-cli/config"
)

// VMDownload
func VMDownload() {

	// download mv
	box, err := os.Create(config.Root + "/nanobox-boot2docker.box")
	if err != nil {
		config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer box.Close()

	//
	FileProgress(fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.box"), box)

	//
	// download vm md5
	md5, err := os.Create(config.Root + "/nanobox-boot2docker.md5")
	if err != nil {
		config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer md5.Close()

	//
	FileDownload("https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.md5", md5)
}
