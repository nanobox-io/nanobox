// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package config

//
import semver "github.com/coreos/go-semver/semver"

//
const VERSION = "0.12.12"

//
var Version *semver.Version

// init
func init() {

	// set the version
	if Version, err = semver.NewVersion(VERSION); err != nil {
		Fatal("[version] semver.NewVersion() failed", err.Error())
	}
}
