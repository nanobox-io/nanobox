//
package server

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/term"

	"github.com/nanobox-io/nanobox/config"
	notifyutil "github.com/nanobox-io/nanobox/util/notify"
	"github.com/nanobox-io/nanobox/util/server/terminal"
)

// Exec
func Exec(params string) error {

	// connect to the server; this will wait until a single read is returned from
	// the server (blocking)
	conn, data, err := connect(fmt.Sprintf("POST /exec?pid=%d&%v HTTP/1.1\r\n\r\n", os.Getpid(), params))
	if err != nil {
		return err
	}
	defer conn.Close()

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

	// print the first read data from above
	os.Stderr.WriteString(string(data))

	//
	return pipeToConnection(conn, stdIn, stdOut)
}
