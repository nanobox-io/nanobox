// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
// +build windows

//
package server

import (
	"github.com/nanobox-io/nanobox/util/server/terminal"
	"os"
	"time"
)

func monitorTerminal(stdOutFD uintptr, params string) {
	tick := time.Tick(time.Millisecond * 250)

	//
	prevW, prevH := terminal.GetTTYSize(stdOutFD)
	for {
		select {
		case <-tick:
			w, h := terminal.GetTTYSize(stdOutFD)

			if prevW != w || prevH != h {
				resizeTTY(stdOutFD, params, w, h)
			}

			prevW, prevH = w, h
		}
	}
}
