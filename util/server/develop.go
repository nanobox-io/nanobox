//
package server

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/term"
	mistClient "github.com/nanopack/mist/core"

	"github.com/nanobox-io/nanobox/config"
	notifyutil "github.com/nanobox-io/nanobox/util/notify"
	"github.com/nanobox-io/nanobox/util/server/terminal"
)

// Develop
func Develop(params string, mist mistClient.Client) error {

	// connect to the server; this will wait until a single read is returned from
	// the server (blocking)
	conn, data, err := connect(fmt.Sprintf("POST /develop?pid=%d&%v HTTP/1.1\r\n\r\n", os.Getpid(), params))
	if err != nil {
		return err
	}
	defer conn.Close()

	// disconnect mist
	mist.Close()

	// begin watching for changes to the project
	go func() {
		if err := notifyutil.Watch(config.CWDir, NotifyServer); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	//
	os.Stderr.WriteString(fmt.Sprintf(`+> Opening a nanobox console:


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

--------------------------------------------------------------------------------
+ You are in a virtual machine (vm)
+ Your local source code has been mounted into the vm, and changes in either
the vm or local will be mirrored.
+ If you run a server, access it at >> %s
--------------------------------------------------------------------------------
`, config.Nanofile.Domain))

	// get current term info
	stdIn, stdOut, _ := term.StdStreams()

	// connect a raw terminal; if no error is returned defer resetting the terminal
	state, err := terminal.Connect(stdIn, stdOut)
	if err == nil {
		defer terminal.Disconnect(stdIn, state)
	}

	// print the first read data from above
	os.Stderr.Write(data)

	//
	return pipeToConnection(conn, stdIn, stdOut)
}
