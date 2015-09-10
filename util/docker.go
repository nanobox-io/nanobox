// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	ossignal "os/signal"
	"runtime"
	"time"

	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"

	"github.com/pagodabox/nanobox-cli/config"
	// "github.com/pagodabox/nanobox-golang-stylish"
)

// run satisfies the Command interface
type Docker struct {
	Params string
}

// Run
func (d Docker) Run() {

	//
	conn, err := net.Dial("tcp4", config.ServerURI)
	if err != nil {
		fmt.Println(err)
		return
	}

	// forward all the signals to the nanobox server
	forwardAllSignals()

	// fake a web request
	conn.Write([]byte(fmt.Sprintf("POST /exec?%v HTTP/1.1\r\n\r\n", d.Params)))

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
	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)
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
			req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/killexec?signal=%s", config.ServerURI, sig), nil)
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

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/resizeexec?h=%d&w=%d", config.ServerURI, h, w), nil)

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
