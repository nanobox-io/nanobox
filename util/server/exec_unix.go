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
	"fmt"
	"os"

	"github.com/docker/docker/pkg/signal"
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
	stdIn, stdOut, _ := term.StdStreams()
	stdInFD, _ := term.GetFdInfo(stdIn)
	stdOutFD, _ := term.GetFdInfo(stdOut)
	// stdErrFD, _ := term.GetFdInfo(stdErr)

	// handle all incoming os signals and act accordingly; default behavior is to
	// forward all signals to nanobox server
	sigs := make(chan os.Signal, 1)
	go func() {
		select {

		// resize the raw terminal w/e the local terminal is resized
		case sig := <-sigs:
			if sig == signal.SIGWINCH {
				resizeTTY(stdOutFD, params)
			}
		}
	}()

	// setup a raw terminal
	oldState, err := term.SetRawTerminal(stdInFD)
	if err != nil {
		config.Fatal("[util/server/exec_unix] term.SetRawTerminal() failed - ", err.Error())
	}
	defer term.RestoreTerminal(stdOutFD, oldState)

	switch where {
	case "exec", "console":
		where = fmt.Sprintf("POST /exec?%v HTTP/1.1\r\n\r\n", params)
	case "develop":
		where = fmt.Sprintf("POST /develop?%v HTTP/1.1\r\n\r\n", params)
	}

	//
	return WriteTCP(where)
}

// resizeTTY
func resizeTTY(fd uintptr, params string) {

	// get current size
	w, h := terminal.GetTTYSize(fd)

	//
	res, err := Post(fmt.Sprintf("/resizeexec?w=%d&h=%d&%v", w, h, params), "text/plain", nil)
	if err != nil {
		fmt.Printf("Error issuing resize: %s\n", err)
	}
	defer res.Body.Close()
}
