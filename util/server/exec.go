//
package server

import (
	"fmt"

	"github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox/config"
	notifyutil "github.com/nanobox-io/nanobox/util/notify"
	"github.com/nanobox-io/nanobox/util/server/terminal"
)

// Exec
func Exec(params string) error {

	// connect to the server
	conn, err := connect(params)
	if err != nil {
		return err
	}

	// begin watching for changes to the project
	go func() {
		if err := notifyutil.Watch(config.CWDir, NotifyServer); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	// get current term info
	stdIn, stdOut, _ := term.StdStreams()

	//
	terminal.Connect(stdIn, stdOut)

	//
	return pipeToConnection(conn, stdIn, stdOut)
}
