// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package config

//
import (
	"fmt"
	"os"

	semver "github.com/coreos/go-semver/semver"
)

//
const VERSION = "0.12.6"

//
var Version *semver.Version

// init
func init() {

	// set the version
	if Version, err = semver.NewVersion(VERSION); err != nil {
		fmt.Println("Fatal error! See ~/.nanobox/nanobox.log for details. Exiting...")
		Log.Fatal("[version] semver.NewVersion() failed", err)
		Log.Close()
		os.Exit(1)
	}
}
