package processors

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/nanobox-io/nanobox/util/odin"
)

// Tunnel ...
func Tunnel(tunnelConfig TunnelConfig) error {

	// clean the app id
	tunnelConfig.App = getAppID(tunnelConfig.App)

	var err error
	var port int
	key, location, container, port, err = odin.EstablishTunnel(getAppID(tunnelConfig.App), tunnelConfig.Container)
	if err != nil {
		return err
	}

	if tunnelConfig.Port == "" {
		tunnelConfig.Port = fmt.Sprintf("%d", port)
	}

	// establish a connection and just leave it open.
	req, err := http.NewRequest("POST", fmt.Sprintf("/tunnel?key=%s", key), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// set noproxy because this connection allows more multiple connections
	// to use the tunnel
	req.Header.Set("X-NOPROXY", "true")
	conn, err := connect(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	//
	serv, err := net.Listen("tcp4", fmt.Sprintf(":%s", tunnelConfig.Port))
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("listening on port", tunnelConfig.Port)

	//
	for {
		conn, err := serv.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
	fmt.Println("done")
	return nil
}

// handleConnection ...
func handleConnection(conn net.Conn) {
	fmt.Println("handling new connection")
	req, err := http.NewRequest("POST", fmt.Sprintf("/tunnel?key=%s", key), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	remoteConn, err := connect(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer remoteConn.Close()

	go io.Copy(conn, remoteConn)
	_, err = io.Copy(remoteConn, conn)
	if err != nil {
		fmt.Println(err)
	}
}
