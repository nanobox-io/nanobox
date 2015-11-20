//
package server

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox/config"
	notifyutil "github.com/nanobox-io/nanobox/util/notify"
	"github.com/nanobox-io/nanobox/util/server/terminal"
)

// Exec
func Exec(where, params string) error {

	//
	stdIn, stdOut, _ := term.StdStreams()

	//
	return execInternal(where, params, stdIn, stdOut)
}

func execInternal(where, params string, in io.Reader, out io.Writer) error {

	// if we can't connect to the server, lets bail out early
	conn, err := net.Dial("tcp4", config.ServerURI)
	if err != nil {
		return err
	}
	defer conn.Close()

	// get current term info
	stdInFD, isTerminal := term.GetFdInfo(in)
	stdOutFD, _ := term.GetFdInfo(out)

	//
	terminal.PrintNanoboxHeader(where)

	// begin watching for changes to the project
	go func() {
		if err := notifyutil.Watch(config.CWDir, NotifyServer); err != nil {
			fmt.Printf(err.Error())
		}
	}()


	// if we are using a term, lets upgrade it to RawMode
	if isTerminal {

		// handle all incoming os signals and act accordingly; default behavior is to
		// forward all signals to nanobox server
		go monitorTerminal(stdOutFD)
		
		oldState, err := term.SetRawTerminal(stdInFD)
		// we only use raw mode if it is available.
		if err == nil {
			defer term.RestoreTerminal(stdInFD, oldState)
		}
	}

	// make a http request
	switch where {
	case "develop":
		if _, err := fmt.Fprintf(conn, "POST /develop?pid=%d&%v HTTP/1.1\r\n\r\n", os.Getpid(), params); err != nil {
			return err
		}
	default:
		if _, err := fmt.Fprintf(conn, "POST /exec?pid=%d&%v HTTP/1.1\r\n\r\n", os.Getpid(), params); err != nil {
			return err
		}
	}

	return pipeToConnection(conn, in, out)
}

// IsContainerExec
func IsContainerExec(args []string) bool {

	// fetch services to see if the command is trying to run on a specific container
	var services []Service
	res, err := Get("/services", &services)
	if err != nil {
		Fatal("[util/server/exec] Get() failed", err.Error())
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
		Fatal("[util/server/exec_unix] Post() failed", err.Error())
	}
	res.Body.Close()
}

// resizeTTY
func resizeTTY(w, h int) {

	//
	res, err := Post(fmt.Sprintf("/resizeexec?pid=%d&w=%d&h=%d", os.Getpid(), w, h), "text/plain", nil)
	if err != nil {
		fmt.Printf("Error issuing resize: %s\n", err)
	}
	res.Body.Close()
}
