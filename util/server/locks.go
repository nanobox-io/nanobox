//
package server

import (
	"fmt"
	"github.com/nanobox-io/nanobox/config"
	"net"
	"time"
)

var lock net.Conn //

// Lock opens a 'lock' with the server; this is done so that nanobox can know how
// many clients are currenctly connected to the server
func Lock() {
	conn, err := net.Dial("tcp", config.ServerURI)
	if err != nil {
		config.Fatal("[util/server/locks] net.Dial() failed", err.Error())
	}

	conn.Write([]byte(fmt.Sprintf("PUT /lock? HTTP/1.1\r\n\r\n")))

	//
	lock = conn
}

// Unlock closes the suspended 'lock' connection to the server indicating that a
// client has disconnected
func Unlock() {

	// if the lock (conn) hasn't already been closed, close it
	if lock != nil {
		lock.Close()
	}

	// this sleep is important because there needs to be enough time for the guest
	// machine to register that our connection has been broken, before we ask if
	// the machine can be suspended (w/o this there is a race condition)
	time.Sleep(1 * time.Second)
}
