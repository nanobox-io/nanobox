// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

//
package server

import (
	"fmt"
	"github.com/docker/docker/pkg/term"
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/notify"
	"github.com/nanobox-io/nanobox/util/server/terminal"
	"io"
	"net"
	"os"
)

// Exec
func Exec(where, kind, params string) error {
	stdIn, stdOut, _ := term.StdStreams()
	return execInternal(where, kind, params, stdIn, stdOut)
}

func execInternal(where, kind, params string, in io.Reader, out io.Writer) error {

	// if we can't connect to the server, lets bail out early
	conn, err := net.Dial("tcp4", config.ServerURI)
	if err != nil {
		return err
	}
	defer conn.Close()

	// get current term info
	stdInFD, isTerminal := term.GetFdInfo(in)
	stdOutFD, _ := term.GetFdInfo(out)

	// if it is not a terminal, then we don't print a header. this allows piping
	// to and from the nanobox command
	if isTerminal {
		terminal.PrintNanoboxHeader(kind, out)
	}

	// begin watching for changes to the project
	go func() {
		if err := notify.Watch(config.CWDir, NotifyServer); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	// handle all incoming os signals and act accordingly; default behavior is to
	// forward all signals to nanobox server
	go handleSignals(stdOutFD, params)

	// if we are using a term, lets upgrade it to RawMode
	if isTerminal {
		oldState, err := term.SetRawTerminal(stdInFD)
		if err != nil {
			config.Fatal("[util/server/exec_unix] term.SetRawTerminal() failed - ", err.Error())
		}
		defer term.RestoreTerminal(stdInFD, oldState)
	}

	// make a http request
	switch where {
	case "develop":
		fmt.Println("got a develop")
		_, err := fmt.Fprintf(conn, "POST /develop?%v HTTP/1.1\r\n\r\n", params)
		if err != nil {
			return err
		}
	default:
		fmt.Println("got an exec")
		_, err := fmt.Fprintf(conn, "POST /exec?%v HTTP/1.1\r\n\r\n", params)
		if err != nil {
			return err
		}
	}

	fmt.Println("starting to pipe")

	return pipeToConnection(conn, in, out)
}

// IsContainerExec
func IsContainerExec(args []string) bool {

	// fetch services to see if the command is trying to run on a specific container
	var services []Service
	res, err := Get("/services", &services)
	if err != nil {
		Fatal("[util/server/exec] Get() failed - ", err.Error())
	}
	res.Body.Close()

	// make an exception for build1, as it wont show up on the list, but will always exist
	if args[0] == "build1" {
		return true
	}

	// look for a match
	for _, service := range services {
		if args[0] == service.Name {
			return true
		}
	}
	return false
}

func sendSignal(sig os.Signal) {
	res, err := Post(fmt.Sprintf("/killexec?signal=%v", sig), "text/plain", nil)
	if err != nil {
		Fatal("[util/server/exec_unix] Post() failed - ", err.Error())
	}
	res.Body.Close()
}

// resizeTTY
func resizeTTY(fd uintptr, params string, w, h int) {

	//
	res, err := Post(fmt.Sprintf("/resizeexec?w=%d&h=%d&%v", w, h, params), "text/plain", nil)
	if err != nil {
		fmt.Printf("Error issuing resize: %s\n", err)
	}
	res.Body.Close()
}
