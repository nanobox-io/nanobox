package processor

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/nanobox-io/nanobox/util/productionAPI"
)

// processTunnel ...
type processTunnel struct {
	control ProcessControl
}

//
func init() {
	Register("tunnel", tunnelFunc)
}

//
func tunnelFunc(control ProcessControl) (Processor, error) {
	return processTunnel{control}, nil
}

//
func (tunnel processTunnel) Results() ProcessControl {
	return tunnel.control
}

//
func (tunnel processTunnel) Process() error {
	var err error
	app := getAppID(tunnel.control.Meta["alias"])
	key, location, container, err = productionAPI.EstablishTunnel(app, tunnel.control.Meta["container"])
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

	serv, err := net.Listen("tcp4", fmt.Sprintf(":%v", tunnel.control.Meta["port"]))
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		conn, err := serv.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
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
