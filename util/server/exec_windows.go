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
	syscall "github.com/docker/docker/pkg/signal"
	"github.com/nanobox-io/nanobox/util/server/terminal"
	"os"
	"time"
)

func handleSignals(stdOutFD uintptr, params string) {
	go monitorAndResizeTTY(fd, params)

	sigs := make(chan os.Signal, 1)
	syscall.CatchAll(sigs)
	defer syscall.StopCatch(sigs)

	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGCHLD: // ignore this one
		// tell nanobox server about the signal
		default:
			sendSignal(sig)
		}
	}
}

// monitorAndResizeTTY
func monitorAndResizeTTY(fd uintptr, params string) {
	tick := time.Tick(time.Millisecond * 250)

	//
	prevW, prevH := terminal.GetTTYSize(fd)
	for {
		select {
		case <-tick:
			w, h := terminal.GetTTYSize(fd)

			if prevW != w || prevH != h {
				resizeTTY(fd, params, w, h)
			}

			prevW, prevH = w, h
		}
	}
}
