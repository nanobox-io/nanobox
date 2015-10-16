// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package print

type (
	print struct{}
	Print interface {
		Verbose(msg string)
		Silence(msg string)
		Color(msg string, v ...interface{})
		Prompt(p string, v ...interface{}) string
		Password(p string) string
	}
)

var (
	Default Print = print{}
)

func (print) Verbose(msg string) {
	Verbose(msg)
}

func (print) Silence(msg string) {
	Silence(msg)
}

func (print) Color(msg string, v ...interface{}) {
	Color(msg, v...)
}

func (print) Prompt(p string, v ...interface{}) string {
	return Prompt(p, v)
}

func (print) Password(p string) string {
	return Password(p)
}
