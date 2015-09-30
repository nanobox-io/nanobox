// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/nanobox-io/nanobox-cli/config"
)

//
func MD5sMatch(localPath, remotePath string) bool {

	// get local md5
	f, err := os.Open(localPath)

	// if there is no local md5 return false
	if err != nil {
		return false
	}
	defer f.Close()

	localMD5, err := ioutil.ReadAll(f)
	if err != nil {
		config.Fatal("[commands/boxInstall] ioutil.ReadAll() failed", err.Error())
	}

	// get remote md5
	res, err := http.Get(remotePath)
	if err != nil {
		config.Fatal("[commands/boxInstall] http.Get() failed", err.Error())
	}
	defer res.Body.Close()

	remoteMD5, err := ioutil.ReadAll(res.Body)
	if err != nil {
		config.Fatal("[commands/boxInstall] ioutil.ReadAll() failed", err.Error())
	}

	return string(localMD5) == string(remoteMD5)
}
