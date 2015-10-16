// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package util

import (
	"io/ioutil"
	"net/http"
	"os"
)

type (
	util struct{}
	Util interface {
		MD5sMatch(string, string) (bool, error)
	}
)

var (
	Default Util = util{}
)

func (_ util) MD5sMatch(localPath, remotePath string) (bool, error) {
	return MD5sMatch(localPath, remotePath)
}

// MD5sMatch determines if a local MD5 matches a remote MD5
func MD5sMatch(localPath, remotePath string) (bool, error) {

	// get local md5
	f, err := os.Open(localPath)

	// if there is no local md5 return false
	if err != nil {
		return false, nil
	}
	defer f.Close()

	localMD5, err := ioutil.ReadAll(f)
	if err != nil {
		return false, err
	}

	// get remote md5
	res, err := http.Get(remotePath)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	remoteMD5, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	return string(localMD5) == string(remoteMD5), nil
}
