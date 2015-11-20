// +build windows

package server

import (
	"time"

	"github.com/nanobox-io/nanobox/util/server/terminal"
)

func monitorTerminal(stdOutFD uintptr) {
	tick := time.Tick(time.Millisecond * 250)

	//
	prevW, prevH := terminal.GetTTYSize(stdOutFD)
	// inform the server what the starting size is
	resizeTTY(prevW, prevH)

	for {
		select {
		case <-tick:
			w, h := terminal.GetTTYSize(stdOutFD)

			if prevW != w || prevH != h {
				resizeTTY(w, h)
			}

			prevW, prevH = w, h
		}
	}
}
