package nanoagent

import (
  "fmt"
  "crypto/tls"
  "net"
  "net/http"
)

func connect(req *http.Request, location string) (net.Conn, error) {

	// if we can't connect to the server, lets bail out early
	conn, err := tls.Dial("tcp4", location, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return conn, fmt.Errorf("failed to establish connection to nanoagent: %s", err.Error())
	}

	// we dont defer a conn.Close() here because we're returning the conn and
	// want it to remain open

	// make an http request
	if err := req.Write(conn); err != nil {
		return conn, fmt.Errorf("failed to establish console session with nanoagent: %s", err.Error())
	}

	return conn, nil
}
