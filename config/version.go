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
const VERSION = "0.9.1"

//
var Version *semver.Version

// init
func init() {

	// set the version; if this fails just panic (since it will be very rare)
	if Version, err = semver.NewVersion(VERSION); err != nil {
		panic(err)
	}
}
