// +build windows

package server

import (
	"github.com/nanobox-io/nanobox/util/server/terminal"
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
