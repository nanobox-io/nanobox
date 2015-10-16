// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package util

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
