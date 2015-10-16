// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package config

//
var (
	Default Config = config{}
)

type (
	Config interface {
		Fatal(string, string)
		Root() string
	}

	config struct {
	}
)

func (_ config) Fatal(msg, err string) {
	Fatal(msg, err)
}

func (_ config) Root() string {
	return Root
}
