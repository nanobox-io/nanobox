package processors

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/nanobox-io/nanobox/util/odin"
)

// Tunnel ...
type Tunnel struct {
	App       string
	Port      string
	Container string
}

//
func (tunnel Tunnel) Run() error {

	// clean the app id
	tunnel.App = getAppID(tunnel.App)

	var err error

	key, location, container, err = odin.EstablishTunnel(getAppID(tunnel.App), tunnel.Container)
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
	conn, err := connect(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()
	fmt.Println("server connection established")

	//
	serv, err := net.Listen("tcp4", fmt.Sprintf(":%v", tunnel.Port))
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("listening on port", tunnel.Port)

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
	fmt.Println("new connectin to server")
	defer remoteConn.Close()

	fmt.Println("piping")
	go io.Copy(conn, remoteConn)
	_, err = io.Copy(remoteConn, conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("connection closed")
}
