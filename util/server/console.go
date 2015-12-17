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

// Console
func Console(params string) error {

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

	//
	os.Stderr.WriteString(`+> Opening a nanobox console:


                                   **
                                ********
                             ***************
                          *********************
                            *****************
                          ::    *********    ::
                             ::    ***    ::
                           ++   :::   :::   ++
                              ++   :::   ++
                                 ++   ++
                                    +

                    _  _ ____ _  _ ____ ___  ____ _  _
                    |\ | |__| |\ | |  | |__) |  |  \/
                    | \| |  | | \| |__| |__) |__| _/\_
`)

	// get current term info
	stdIn, stdOut, _ := term.StdStreams()

	//
	terminal.Connect(stdIn, stdOut)

	//
	return pipeToConnection(conn, stdIn, stdOut)
}
