// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"
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
// NOTE: a sibling command to this command is 'nanobox console', which simply
// calls this requesting a console rather than execing a command. The commands are
// basically aliases for eachother from the servers perspective, however from the
// users perspective it's all about context.
//
// In an effort to avoid duplicating a lot code one command just calls the other.
// Initially exec was calling console, which would allow for each command to
// gracefully failover to the other. While this is a good thing, we intentially
// chose to hardlock the functionality of each command to avoid any confusion.
//
// Since exec is the command that has additional options, it's going to be the
// one that does all the checking, and decide if the command will be a console
// or exec, providing us with our hardlocking.
func (c *ExecCommand) Run(opts []string) {

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	//
	var fTunnel string
	flags.StringVar(&fTunnel, "t", "", "")
	flags.StringVar(&fTunnel, "tunnel", "", "")

	//
	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.destroy] flags.Parse() failed", err)
	}

	//
	cmd := flags.Args()[0:]

	//
	if !c.console && len(cmd) <= 0 {
		cmd = append(cmd, ui.Prompt("Please specify a command you wish to exec: "))
	}

	//
	if c.console && len(flags.Args()) > 0 {
		fmt.Println("Attempting to run 'nanobox console' with a command. Did you intend to run 'nanobox exec'?")
		os.Exit(0)
	}

	//
	if len(flags.Args()) >= 1 {

	}

	// add a check here to regex the fTunnel to make sure the format makes sense

	//
	v := url.Values{}
	v.Add("forward", fTunnel)
	v.Add("cmd", strings.Join(cmd, " "))

	//
	connect(fmt.Sprintf("POST /exec?%s HTTP/1.1\r\n\r\n", v.Encode()))
}

//
func connect(path string) {
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
	conn.Write([]byte(path))

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
