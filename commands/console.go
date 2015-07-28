// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	ossignal "os/signal"
	"runtime"
	"time"

	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	// "github.com/pagodabox/nanobox-golang-stylish"
)

// ConsoleCommand satisfies the Command interface
type ConsoleCommand struct{}

// Help
func (c *ConsoleCommand) Help() {
	ui.CPrint(`
Description:
  Drops you into bash inside your nanobox vm

Usage:
  nanobox exec
  `)
}

// Run
func (c *ConsoleCommand) Run(opts []string) {

	conn, err := net.Dial("tcp4", fmt.Sprintf("%v:1757", config.Nanofile.IP))
	if err != nil {
		fmt.Println(err)
		return
	}
	// forward all the signals to the nanobox server
	forwardAllSignals()

	// make sure we dont just kill the connection
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	v := url.Values{}

	// if there are any options, we should expect opts[0] to be a command that is
	// wanting to be run
	if len(opts) >= 1 {
		v.Add("cmd", opts[0])
	}

	// fake a web request
	conn.Write([]byte(fmt.Sprintf("POST /exec?%v HTTP/1.1\r\n\r\n", v.Encode())))

	// setup a raw terminal
	var oldState *term.State
	stdIn, stdOut, _ := term.StdStreams()
	inFd, _ := term.GetFdInfo(stdIn)
	outFd, _ := term.GetFdInfo(stdOut)

	// monitor the window size and send a request whenever we resize
	monitorSize(outFd)

	oldState, err = term.SetRawTerminal(inFd)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.RestoreTerminal(inFd, oldState)

	// pipe data
	go io.Copy(os.Stdout, conn)
	io.Copy(conn, os.Stdin)
}

// forwardAllSignals forwards the signals you recieve and send them to nanobox server
func forwardAllSignals() {
	sigc := make(chan os.Signal, 128)
	signal.CatchAll(sigc)
	go func() {
		for s := range sigc {
			// skip children and window resizes
			if s == signal.SIGCHLD || s == signal.SIGWINCH {
				continue
			}
			var sig string
			for sigStr, sigN := range signal.SignalMap {
				if sigN == s {
					sig = sigStr
					break
				}
			}

			//
			if sig == "" {
				fmt.Printf("Unsupported signal: %v. Discarding.\n", s)
			}

			//
			req, _ := http.NewRequest("POST", fmt.Sprintf("http://nanobox-ruby-sample.nano.dev:1757/killexec?signal=%s", sig), nil)
			_, err := http.DefaultClient.Do(req)
			fmt.Println(err)
		}
	}()
	return
}

// monitorSize
func monitorSize(outFd uintptr) {
	resizeTTY(outFd)

	if runtime.GOOS == "windows" {
		go func() {
			prevH, prevW := getTTYSize(outFd)
			for {
				time.Sleep(time.Millisecond * 250)
				h, w := getTTYSize(outFd)

				if prevW != w || prevH != h {
					resizeTTY(outFd)
				}
				prevH = h
				prevW = w
			}
		}()
	} else {
		sigchan := make(chan os.Signal, 1)
		ossignal.Notify(sigchan, signal.SIGWINCH)
		go func() {
			for range sigchan {
				resizeTTY(outFd)
			}
		}()
	}
}

// resizeTTY
func resizeTTY(outFd uintptr) {
	h, w := getTTYSize(outFd)

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://nanobox-ruby-sample.nano.dev:1757/resizeexec?h=%d&w=%d", h, w), nil)

	//
	if _, err := http.DefaultClient.Do(req); err != nil {
		fmt.Println("ERR!", err)
	}
}

// getTTYSize
func getTTYSize(outFd uintptr) (h, w int) {
	ws, err := term.GetWinsize(outFd)
	if err != nil {
		fmt.Printf("Error getting size: %s\n", err)
		if ws == nil {
			return 0, 0
		}
	}
	h = int(ws.Height)
	w = int(ws.Width)
	return
}
