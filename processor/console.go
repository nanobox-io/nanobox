// +build !windows

package processor

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	syscall "github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox/util/odin"
)

// processConsole ...
type processConsole struct {
	control   ProcessControl
	container string
	app       string
}

//
func init() {
	Register("console", consoleFn)
}

//
func consoleFn(control ProcessControl) (Processor, error) {
	console := &processConsole{control: control}
	return console, console.validateMeta()
}

//
func (console processConsole) Results() ProcessControl {
	return console.control
}

//
func (console processConsole) Process() error {

	var err error

	// get a key from odin
	key, location, container, err = odin.EstablishConsole(getAppID(console.app), console.container)
	if err != nil {
		return err
	}
	fmt.Println("key", key, "location", location, "container", container)
	// establish connection to nanoagent
	req, err := http.NewRequest("POST", fmt.Sprintf("/exec?key=%s&container=%s", key, container), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// connect to remote machine
	remoteConn, bytes, err := connect(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer remoteConn.Close()
	stdIn, stdOut, _ := term.StdStreams()

	// establish file descriptors
	stdInFD, isTerminal := term.GetFdInfo(stdIn)
	stdOutFD, _ := term.GetFdInfo(stdOut)

	// if we are using a term, lets upgrade it to RawMode
	if isTerminal {

		// handle all incoming os signals and act accordingly; default behavior is to
		// forward all signals to nanobox server
		go monitor(stdOutFD)

		oldState, err := term.SetRawTerminal(stdInFD)

		// we only use raw mode if it is available.
		if err == nil {
			defer term.RestoreTerminal(stdInFD, oldState)
		}

		// pipe data from out to the conn (server)
		go func() {
			io.Copy(remoteConn, stdIn)
		}()

		os.Stdout.Write(bytes)

		io.Copy(stdOut, remoteConn)
		remoteConn.(*tls.Conn).Close()
	}
	return nil
}

// validateMeta validates that the required metadata exists
func (console *processConsole) validateMeta() error {

	// set optional meta values
	console.app = console.control.Meta["app"]

	// set container (required) and ensure it's provided
	console.container = console.control.Meta["container"]
	if console.container == "" {
		return fmt.Errorf("Missing required meta value 'container'")
	}

	return nil
}

// monitor ...
func monitor(stdOutFD uintptr) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGWINCH)
	defer signal.Stop(sigs)

	// inform the server what the starting size is
	resizeTTY(getTTYSize(stdOutFD))

	// resize the tty for any signals received
	for range sigs {
		resizeTTY(getTTYSize(stdOutFD))
	}
}

// getTTYSize ...
func getTTYSize(fd uintptr) (int, int) {

	ws, err := term.GetWinsize(fd)
	if err != nil {
		fmt.Printf("GetWinsize() failed - %s\n", err.Error())
		os.Exit(1)
	}

	//
	return int(ws.Width), int(ws.Height)
}

// resizeTTY ...
func resizeTTY(w, h int) {
	//
	if _, err := http.Post(fmt.Sprintf("https://%s/resizeexec?key=%s&container=%s&w=%d&h=%d", location, key, container, w, h), "plain/text", nil); err != nil {
		fmt.Printf("Error issuing resize: %s\n", err)
	}
}
