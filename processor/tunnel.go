package processor

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/nanobox-io/nanobox/util/production_api"
)

type tunnel struct {
	config ProcessConfig
}

func init() {
	Register("tunnel", tunnelFunc)
}

func tunnelFunc(config ProcessConfig) (Processor, error) {
	return tunnel{config}, nil
}

func (self tunnel) Results() ProcessConfig {
	return self.config
}

func (self tunnel) Process() error {
	var err error
	app := getAppID(self.config.Meta["alias"])
	key, location, container, err = production_api.EstablishTunnel(app, self.config.Meta["container"])
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

	serv, err := net.Listen("tcp4", ":"+self.config.Meta["port"])
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
