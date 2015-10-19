// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

//
package server

import (
	"errors"
	"fmt"
	terminal "github.com/docker/docker/pkg/term"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/notify"
	"io"
	"net"
	"os"
	"time"
)

var (
	DisconnectedFromServer = errors.New("the server went away")
)

// Exec
func Exec(kind, params string) error {
	stdIn, stdOut, _ := terminal.StdStreams()
	return execInternal(kind, params, stdIn, stdOut)
}

func execInternal(kind, params string, in io.Reader, out io.Writer) error {

	// if we can't connect to the server, lets bail out early
	conn, err := net.Dial("tcp4", config.ServerURI)
	if err != nil {
		return err
	}
	defer conn.Close()

	printNanoboxHeader(kind)

	// begin watching for changes to the project
	go func() {
		if err := notify.Watch(config.CWDir, NotifyServer); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	// get current terminal info
	stdInFD, isTerminal := terminal.GetFdInfo(in)
	stdOutFD, _ := terminal.GetFdInfo(out)

	// handle all incoming os signals and act accordingly; default behavior is to
	// forward all signals to nanobox server
	go handleSignals(stdOutFD, params)

	// if we are using a terminal, lets upgrade it to RawMode
	if isTerminal {
		oldState, err := terminal.SetRawTerminal(stdInFD)
		if err != nil {
			config.Fatal("[util/server/exec_unix] terminal.SetRawTerminal() failed - ", err.Error())
		}
		defer terminal.RestoreTerminal(stdInFD, oldState)
	}

	// set up notification channels so that when a ping check fails, we can disconnect the active console
	disconnect := make(chan error)
	done := make(chan interface{})
	go monitorServer(done, disconnect, 5*time.Second)

	// create a web request that will be upgrated to a raw socket on the server
	conn.Write([]byte(fmt.Sprintf("POST /exec?%v HTTP/1.1\r\n\r\n", params)))

	// pipe data from the server to out, and from in to the server
	go func() {
		go io.Copy(os.Stdout, conn)
		io.Copy(conn, os.Stdin)
		close(done)
	}()
	return <-disconnect
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

//
func printNanoboxHeader(kind string) {
	switch kind {

	//
	case "command":
		fmt.Printf(stylish.Bullet("Executing command in nanobox..."))

		//
	case "console", "container":
		fmt.Printf(`+> Opening a nanobox console:


                                     **
                                  ********
                               ***************
                            *********************
                              *****************
                            ::    *********    ::
                               ::    ***    ::
                             ++   :::   :::   ++
                                ++   :::   ++
                                   ++   ++
                                      +

                      _  _ ____ _  _ ____ ___  ____ _  _
                      |\ | |__| |\ | |  | |__) |  |  \/
                      | \| |  | | \| |__| |__) |__| _/\_
`)

		if kind == "console" {
			fmt.Printf(`
--------------------------------------------------------------------------------
+ You are in a virtual machine (vm)
+ Your local source code has been mounted into the vm, and changes in either
the vm or local will be mirrored.
+ If you run a server, access it at >> %s
--------------------------------------------------------------------------------
`, config.Nanofile.Domain)
		}
	}
}

// monitor the server for disconnects
func monitorServer(done chan interface{}, disconnect chan<- error, after time.Duration) {
	defer close(disconnect)
	ping := make(chan interface{}, 1)
	for {
		// ping the server
		go func() {
			if ok, _ := Ping(); ok {
				ping <- true
			} else {
				close(ping)
			}
		}()
		select {
		case <-ping:
		case <-time.After(after):
			disconnect <- DisconnectedFromServer
			return
		case <-done:
			return
		}
	}
}

// getTTYSize
func getTTYSize(fd uintptr) (int, int) {

	ws, err := terminal.GetWinsize(fd)
	if err != nil {
		config.Fatal("[util/server/exec] terminal.GetWinsize() failed", err.Error())
	}

	//
	return int(ws.Width), int(ws.Height)
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

func sendSignal(sig os.Signal) {
	res, err := Post(fmt.Sprintf("/killexec?signal=%v", sig), "text/plain", nil)
	if err != nil {
		Fatal("[util/server/exec_unix] Post() failed - ", err.Error())
	}
	res.Body.Close()
}
