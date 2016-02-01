//
package terminal

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox/config"
)

// Connect
func Connect(in io.Reader, out io.Writer) (*term.State, error) {

	stdInFD, _ := term.GetFdInfo(in)
	stdOutFD, _ := term.GetFdInfo(out)

	// handle all incoming os signals and act accordingly; default behavior is to
	// forward all signals to nanobox server
	go monitor(stdOutFD)

	// try to upgrade to a raw terminal; if accessed via a terminal this will upgrade
	// with no error, if not an error will be returned
	return term.SetRawTerminal(stdInFD)
}

// Disconnect
func Disconnect(in io.Reader, oldState *term.State) {
	stdInFD, _ := term.GetFdInfo(in)
	term.RestoreTerminal(stdInFD, oldState)
}

// getTTYSize
func getTTYSize(fd uintptr) (int, int) {

	ws, err := term.GetWinsize(fd)
	if err != nil {
		config.Fatal("[util/server/exec] term.GetWinsize() failed", err.Error())
	}

	//
	return int(ws.Width), int(ws.Height)
}

// resizeTTY
func resizeTTY(w, h int) {

	//
	if _, err := http.Post(fmt.Sprintf("%s/resizeexec?pid=%d&w=%d&h=%d", config.ServerURL, os.Getpid(), w, h), "text/plain", nil); err != nil {
		fmt.Printf("Error issuing resize: %s\n", err)
	}
}
