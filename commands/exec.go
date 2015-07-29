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
	"strings"
	"time"

	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	// "github.com/pagodabox/nanobox-golang-stylish"
)

// ExecCommand satisfies the Command interface
type ExecCommand struct{
	console bool
}

// Help
func (c *ExecCommand) Help() {
	ui.CPrint(`
Description:
  Run's a command on an application's service.

Usage:
  nanobox exec <COMMAND>
	nanobox exec -t local:remote,local:remote <COMMAND>

	ex. nanobox exec ls -la

Options:
	-t --tunnel
		Establishes a port forward for each comma delimited local:remote port combo
	`)
}

// Run
//
// NOTE: the sibling command to this command is 'nanobox console', which simply
// calls this command in a way that changes some of the flags it passes. The reason
// for doing it this way, is the commands are basically aliases for eachother from
// the servers perspective, however from the users perspective it's all about context.
//
// In an effort to avoid duplicating a lot of the same code to get both commands
// to run, one just calls the other. Initially exec was calling console, which
// reduced a lot of logical checking, however that would allow for each command
// to gracefully failover to the other. While this is a good thing, we intentially
// chose to hardlock the functionality of each command to avoid any confusion.
//
// Since exec is the command that has additional options, it's going to be the
// one that does all the checking, and decide if the command will be a console
// or exec, providing us with our hard locking.
func (c *ExecCommand) Run(opts []string) {

	//
	v := url.Values{}

	//
	var cmd []string

	//
	switch {

	// a command with tunnel flag but no options
	case len(opts) == 1 && opts[0] == "-t":
		fmt.Println("-t flag set but missing options, exiting...")
		os.Exit(1)

	// all console command
	case c.console:
		switch {

		// console with no tunnel flag
		case c.console && len(opts) <= 0:
			fallthrough

		// console with tunnel flag
		case c.console && len(opts) == 2 && opts[0] == "-t":
			v.Add("forward", opts[1])

		// console with command
		case c.console && len(opts) >= 1 && opts[0] != "-t":
			fmt.Println("Attempting to run 'nanobox console' with a command. Did you intend to run 'nanobox exec'?")
			os.Exit(0)

		// console with command and tunnel flag
		case c.console && len(opts) >= 3 && opts[0] == "-t":
			fmt.Println("Attempting to run 'nanobox console' with a command. Did you intend to run 'nanobox exec'?")
			os.Exit(0)
		}

	// all exec related variations
	case !c.console:
		switch {

		// exec missing a command
		case !c.console && len(opts) <= 0:
			cmd = append(cmd, ui.Prompt("Please specify a command you wish to exec: "))

		// exec with command and no tunnel flag
		case !c.console && len(opts) >= 1 && opts[0] != "-t":
			cmd = opts[0:]

		// exec with tunnel flag but no command
		case !c.console && len(opts) == 2 && opts[0] == "-t":
			v.Add("forward", opts[1])
			cmd = append(cmd, ui.Prompt("Please specify a command you wish to exec: "))

		// exec with command and tunnel flag
		case !c.console && len(opts) >= 3 && opts[0] == "-t":
			v.Add("forward", opts[1])
			cmd = opts[2:]
		}

		// something I didn't think of
		default:
			fmt.Println("unrecognized command... exiting")
			os.Exit(1)
	}

	// set the command to be run (if any)
	v.Add("cmd", strings.Join(cmd, " "))

	//
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
