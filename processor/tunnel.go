package processor

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/nanobox-io/nanobox/util/odin"
)

// processTunnel ...
type processTunnel struct {
	control   ProcessControl
	app       string
	port      string
	container string
}

//
func init() {
	Register("tunnel", tunnelFn)
}

//
func tunnelFn(control ProcessControl) (Processor, error) {
	tunnel := &processTunnel{control: control}
	return tunnel, tunnel.validateMeta()
}

//
func (tunnel processTunnel) Results() ProcessControl {
	return tunnel.control
}

//
func (tunnel processTunnel) Process() error {

	var err error

	key, location, container, err = odin.EstablishTunnel(getAppID(tunnel.app), tunnel.container)
	if err != nil {
		return err
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
	conn, _, err := connect(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	//
	serv, err := net.Listen("tcp4", fmt.Sprintf(":%v", tunnel.port))
	if err != nil {
		fmt.Println(err)
		return err
	}

	//
	for {
		conn, err := serv.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}

// validateMeta validates that the required metadata exists
func (tunnel *processTunnel) validateMeta() error {

	// set optional meta values
	tunnel.app = tunnel.control.Meta["app"]

	// set container (required) and ensure it's provided
	tunnel.container = tunnel.control.Meta["container"]
	if tunnel.container == "" {
		return fmt.Errorf("Missing required meta value 'container'")
	}

	// set port (required) and ensure it's provided
	tunnel.port = tunnel.control.Meta["port"]
	if tunnel.port == "" {
		return fmt.Errorf("Missing required meta value 'port'")
	}

	return nil
}

// handleConnection ...
func handleConnection(conn net.Conn) {
	req, err := http.NewRequest("POST", fmt.Sprintf("/tunnel?key=%s", key), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	remoteConn, bytes, err := connect(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer remoteConn.Close()

	conn.Write(bytes)
	go io.Copy(conn, remoteConn)
	_, err = io.Copy(remoteConn, conn)
	if err != nil {
		fmt.Println(err)
	}
}
