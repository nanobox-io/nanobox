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
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"
	"github.com/go-fsnotify/fsnotify"

	"github.com/nanobox-io/nanobox-cli/config"
	// "github.com/nanobox-io/nanobox-golang-stylish"
)

// run satisfies the Command interface
type Server struct {
	Params string
}

// Run
func (d *Server) Run() {

	// begin watching for changes to the project
	go Watch(config.CWDir, handleFileEvent)

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
	monitorSize(outFd, d.Params)

	//
	oldState, err = term.SetRawTerminal(inFd)
	if err != nil {
		fmt.Println(err)
		return
	}

	//
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
			req, _ := http.NewRequest("POST", fmt.Sprintf("%s/killexec?signal=%s", config.ServerURL, sig), nil)
			_, err := http.DefaultClient.Do(req)
			fmt.Println(err)
		}
	}()
	return
}

// handleFileEvent
func handleFileEvent(event *fsnotify.Event, err error) {

	//
	if err != nil {
		if _, ok := err.(syscall.Errno); ok {
			fmt.Printf(`
! WARNING !
Failed to watch files, max file descriptor limit reached. Nanobox will not
be able to propagate filesystem events to the virtual machine. Consider
increasing your max file descriptor limit to re-enable this functionality.
`)
		}

		// need to log this error, or print it out after or something
		return
	}

	//
	if event.Op != fsnotify.Chmod {
		event.Name = strings.Replace(event.Name, config.CWDir, "", -1)

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/file-change?filename=%s", config.ServerURL, event.Name), nil)

		//
		if _, err := http.DefaultClient.Do(req); err != nil {
			// need to log this error, or print it out after or something
			return
		}
	}
}

// monitorSize
func monitorSize(outFd uintptr, params string) {
	resizeTTY(outFd, params)

	if runtime.GOOS == "windows" {
		go func() {
			prevH, prevW := getTTYSize(outFd)
			for {
				time.Sleep(time.Millisecond * 250)
				h, w := getTTYSize(outFd)

				if prevW != w || prevH != h {
					resizeTTY(outFd, params)
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
				resizeTTY(outFd, params)
			}
		}()
	}
}

// resizeTTY
func resizeTTY(outFd uintptr, params string) {
	h, w := getTTYSize(outFd)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/resizeexec?h=%d&w=%d&%v", config.ServerURL, h, w, params), nil)

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
