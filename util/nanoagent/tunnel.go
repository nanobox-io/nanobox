package nanoagent

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

func Tunnel(key, location, port string) error {
	// establish a connection and just leave it open.
	req, err := http.NewRequest("POST", fmt.Sprintf("/tunnel?key=%s", key), nil)
	if err != nil {
		return fmt.Errorf("failed to generate a request for nanoagent: %s", err.Error())
	}

	// set noproxy because this connection allows more multiple connections
	// to use the tunnel
	req.Header.Set("X-NOPROXY", "true")
	conn, err := connect(req, location)
	if err != nil {
		return err
	}
	defer conn.Close()

	// setup a tcp listener
	serv, err := net.Listen("tcp4", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to setup tcp listener: %s", err.Error())
	}

	// fmt.Println("listening on port", tunnelConfig.Port)

	// handle connections
	for {
		conn, err := serv.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept client connection: %s", err.Error())
		}

		go handleConnection(conn, key, location)
	}

	return nil
}

func handleConnection(conn net.Conn, key, location string) {
	req, err := http.NewRequest("POST", fmt.Sprintf("/tunnel?key=%s", key), nil)
	if err != nil {
		return
	}

	remoteConn, err := connect(req, location)
	if err != nil {
		return
	}
	defer remoteConn.Close()

	go io.Copy(conn, remoteConn)
	_, err = io.Copy(remoteConn, conn)
	if err != nil {
		return
	}
}
