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
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/notify"
	"github.com/nanobox-io/nanobox/util/server/terminal"
)

// Exec
func Exec(where, what, params string) error {

	//
	terminal.Header(what)

	// begin watching for changes to the project
	go func() {
		if err := notify.Watch(config.CWDir, NotifyServer); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	// get current terminal info
	stdIn, stdOut, _ := terminal.StdStreams()
	stdInFD, _ := terminal.GetFdInfo(stdIn)
	stdOutFD, _ := terminal.GetFdInfo(stdOut)
	// stdErrFD, _ := terminal.GetFdInfo(stdErr)

	//
	monitorAndResizeTTY(stdOutFD, params)

	// setup a raw terminal
	oldState, err := terminal.SetRawTerminal(stdInFD)
	if err != nil {
		config.Fatal("[util/server/exec_unix] terminal.SetRawTerminal() failed - ", err.Error())
	}
	defer terminal.RestoreTerminal(stdOutFD, oldState)

	switch where {
	case "exec", "console":
		where = fmt.Sprintf("POST /exec?%v HTTP/1.1\r\n\r\n", params)
	case "develop":
		where = fmt.Sprintf("POST /develop?%v HTTP/1.1\r\n\r\n", params)
	}

	// (blocking)
	return WriteTCP(where)
}

// monitorAndResizeTTY
func monitorAndResizeTTY(fd uintptr, params string) {

	//
	prevW, prevH := getTTYSize(fd)

	//
	go func() {
		for {

			w, h := getTTYSize(fd)

			if prevW != w || prevH != h {
				res, err := Post(fmt.Sprintf("/resizeexec?w=%d&h=%d&%v", w, h, params), "text/plain", nil)
				if err != nil {
					fmt.Printf("Error issuing resize: %s\n", err)
				}
				defer res.Body.Close()
			}

			prevW, prevH = w, h

			time.Sleep(time.Millisecond * 250)
		}
	}()
}
