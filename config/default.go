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
		ParseConfig(path string, v interface{}) error
		Debug(msg string, debug bool)
		Info(msg string, debug bool)
		ParseNanofile() *NanofileConfig
		ParseVMfile() *VMfileConfig
	}

	config struct {
	}
)

func (_ config) ParseVMfile() *VMfileConfig {
	return ParseVMfile()
}

func (_ config) ParseNanofile() *NanofileConfig {
	return ParseNanofile()
}

func (_ config) Debug(msg string, debug bool) {
	Debug(msg, debug)
}

func (_ config) Info(msg string, debug bool) {
	Info(msg, debug)

}
func (_ config) ParseConfig(path string, v interface{}) error {
	return ParseConfig(path, v)
}

func (_ config) Fatal(msg, err string) {
	Fatal(msg, err)
}

func (_ config) Root() string {
	return Root
}
