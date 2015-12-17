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
func Connect(in io.Reader, out io.Writer) {

	stdInFD, isTerminal := term.GetFdInfo(in)
	stdOutFD, _ := term.GetFdInfo(out)

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
	}
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
	if _, err := http.Post(fmt.Sprintf("/resizeexec?pid=%d&w=%d&h=%d", os.Getpid(), w, h), "text/plain", nil); err != nil {
		fmt.Printf("Error issuing resize: %s\n", err)
	}
}
