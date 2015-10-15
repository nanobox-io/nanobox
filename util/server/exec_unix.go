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
	"io"
	"net"
	"os"
	"time"

	syscall "github.com/docker/docker/pkg/signal"
	terminal "github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/notify"
)

// Exec
func Exec(kind, params string) error {

	//
	header(kind)

	// begin watching for changes to the project
	go notify.Watch(config.CWDir, NotifyServer)

	// get current terminal info
	stdIn, stdOut, _ := terminal.StdStreams()
	stdInFD, _ := terminal.GetFdInfo(stdIn)
	stdOutFD, _ := terminal.GetFdInfo(stdOut)
	// stdErrFD, _ := terminal.GetFdInfo(stdErr)

	// handle all incoming os signals and act accordingly; default behavior is to
	// forward all signals to nanobox server
	sigs := make(chan os.Signal, 1)
	go func() {
		select {

		// received a signal
		case sig := <-sigs:

			switch sig {

			// skip the following signals
			case syscall.SIGCHLD:
				// do nothing

			// resize the raw terminal w/e the local terminal is resized
			case syscall.SIGWINCH:
				resizeTTY(stdOutFD, params)

			// tell nanobox server about the signal; this kills the terminal
			default:
				res, err := Post(fmt.Sprintf("/killexec?signal=%v", sig), "text/plain", nil)
				if err != nil {
					Fatal("[util/server/exec_unix] Post() failed - ", err.Error())
				}
				defer res.Body.Close()
			}
		}
	}()

	// setup a raw terminal
	oldState, err := terminal.SetRawTerminal(stdInFD)
	if err != nil {
		config.Fatal("[util/server/exec_unix] terminal.SetRawTerminal() failed - ", err.Error())
	}
	defer terminal.RestoreTerminal(stdOutFD, oldState)

	// set up a ping to detect when the server goes away to be able to disconnect
	// any active consoles
	disconnect := make(chan struct{})
	go func() {
		for {
			select {
			default:
				if ok, _ := Ping(); !ok {
					close(disconnect)
				}
				time.Sleep(2 * time.Second)

			// return after timeout
			case <-time.After(5 * time.Second):
				return

			// return if ping fails
			case <-disconnect:
				return
			}
		}
	}()

	//
	conn, err := net.Dial("tcp4", config.ServerURI)
	if err != nil {
		return err
	}

	// pseudo a web request
	conn.Write([]byte(fmt.Sprintf("POST /exec?%v HTTP/1.1\r\n\r\n", params)))

	// pipe data
	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)

	return nil
}

// resizeTTY
func resizeTTY(fd uintptr, params string) {

	// get current size
	w, h := getTTYSize(fd)

	//
	res, err := Post(fmt.Sprintf("/resizeexec?w=%d&h=%d&%v", w, h, params), "text/plain", nil)
	if err != nil {
		fmt.Printf("Error issuing resize: %s\n", err)
	}
	defer res.Body.Close()
}
