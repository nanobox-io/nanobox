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
		Info(msg string)
		Error(msg string)
		ParseNanofile() *NanofileConfig
		ParseVMfile() *VMfileConfig
	}

	config struct {
	}
)

func (config) ParseVMfile() *VMfileConfig {
	return ParseVMfile()
}

func (config) ParseNanofile() *NanofileConfig {
	return ParseNanofile()
}

func (config) Debug(msg string, debug bool) {
	Debug(msg, debug)
}

func (config) Info(msg string) {
	Info(msg)
}

func (config) Error(msg string) {
	Error(msg)
}

func (config) ParseConfig(path string, v interface{}) error {
	return ParseConfig(path, v)
}

func (config) Fatal(msg, err string) {
	Fatal(msg, err)
}

func (config) Root() string {
	return Root
}
