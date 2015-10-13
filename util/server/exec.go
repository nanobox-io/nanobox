// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/notify"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

// Exec
func Exec(kind, params string) error {

	//
	header(kind)

	// begin watching for changes to the project
	go notify.Watch(config.CWDir, NotifyServer)

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
			case <-time.After(5 * time.Second):
				os.Exit(0)
			case <-disconnect:
				os.Exit(0)
			}
		}
	}()

	//
	conn, err := net.Dial("tcp4", config.ServerURI)
	if err != nil {
		return err
	}

	// forward all the signals to the nanobox server
	forwardAllSignals()

	// fake a web request
	conn.Write([]byte(fmt.Sprintf("POST /exec?%v HTTP/1.1\r\n\r\n", params)))

	// setup a raw terminal
	var oldState *term.State
	stdIn, stdOut, _ := term.StdStreams()
	inFd, _ := term.GetFdInfo(stdIn)
	outFd, _ := term.GetFdInfo(stdOut)

	// monitor the window size and send a request whenever we resize
	monitorSize(outFd, params)

	//
	oldState, err = term.SetRawTerminal(inFd)
	if err != nil {
		return err
	}

	//
	defer term.RestoreTerminal(inFd, oldState)

	// pipe data
	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)

	return nil
}

// IsContainerExec
func IsContainerExec(args []string) (found bool) {

	// fetch services to see if the command is trying to run on a specific container
	var services []Service
	res, err := Get("/services", &services)
	if err != nil {
		config.Fatal("[commands/nanoboxExec] failed - ", err.Error())
	}
	defer res.Body.Close()

	//
	for _, service := range services {

		// range over the services to find a potential match for args[0] or make an
		// exception for 'build1' since that wont show up on the list.
		if args[0] == service.Name || args[0] == "build1" {
			found = true
		}
	}

	return
}

//
func header(kind string) {
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
