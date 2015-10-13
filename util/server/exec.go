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
	"runtime"
	"time"

	syscall "github.com/docker/docker/pkg/signal"
	terminal "github.com/docker/docker/pkg/term"

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

			// resize the raw terminal w/e the local terminal is resized
			case syscall.SIGWINCH:
				resizeTTY(stdOutFD, params)

			// skip the following signals
			case syscall.SIGCHLD:
				// do nothing

			// tell nanobox server about the signal; this kills the terminal
			default:
				res, err := Post(fmt.Sprintf("/killexec?signal=%v", sig), "text/plain", nil)
				if err != nil {
					fmt.Println(err)
				}
				defer res.Body.Close()
			}
		}
	}()

	// setup a raw terminal
	oldState, err := terminal.SetRawTerminal(stdInFD)
	if err != nil {
		return err
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

// resizeTTY
func resizeTTY(fd uintptr, params string) error {

	// get current size
	w, h := getTTYSize(fd)

	//
	if runtime.GOOS == "windows" {
		prevW, prevH := getTTYSize(fd)

		for {
			time.Sleep(time.Millisecond * 250)
			w, h := getTTYSize(fd)

			if prevW != w || prevH != h {
				resizeTTY(fd, params)
			}

			prevW, prevH = w, h
		}
	}

	res, err := Post(fmt.Sprintf("/resizeexec?w=%d&h=%d&%v", w, h, params), "text/plain", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// getTTYSize
func getTTYSize(fd uintptr) (w, h int) {

	ws, err := terminal.GetWinsize(fd)
	if err != nil {
		fmt.Printf("Error getting size: %s\n", err)
	}

	//
	if ws == nil {
		w, h = 0, 0
	}

	//
	w, h = int(ws.Width), int(ws.Height)

	return
}
