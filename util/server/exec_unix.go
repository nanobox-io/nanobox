// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
// +build !windows

//
package server

import (
	syscall "github.com/docker/docker/pkg/signal"
	"github.com/nanobox-io/nanobox/util/server/terminal"
	"os"
)

func handleSignals(stdOutFD uintptr, params string) {
	sigs := make(chan os.Signal, 1)
	syscall.CatchAll(sigs)
	defer syscall.StopCatch(sigs)
	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGCHLD: // ignore this one
		case syscall.SIGWINCH:
			w, h := terminal.GetTTYSize(stdOutFD)
			resizeTTY(stdOutFD, params, w, h)
		// tell nanobox server about the signal
		default:
			sendSignal(sig)
		}
	}
}
