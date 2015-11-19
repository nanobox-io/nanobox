//
package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/nanobox-io/nanobox/config"
)

var (
	DisconnectedFromServer = errors.New("the server went away")
)

// Service
type Service struct {
	CreatedAt time.Time
	EnvVars   map[string]string
	IP        string
	Name      string
	Password  string
	Ports     []int
	Username  string
	UID       string
}

// Get issues a "GET" request to nanobox server at the specified path
func Get(path string, v interface{}) (*http.Response, error) {

	req, err := http.NewRequest("GET", config.ServerURL+path, nil)
	if err != nil {
		return nil, err
	}

	//
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read the body
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal response into the provided container. If no container was given
	// it's mostly likely a raw request where the body isn't needed, and therfore
	// this step can be skipped.
	if v != nil {
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// Post issues a "POST" to nanobox server
func Post(path, contentType string, body io.Reader) (*http.Response, error) {

	res, err := http.Post(config.ServerURL+path, contentType, body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Put
func Put(path string, body io.Reader) (*http.Response, error) {

	// check to see if the VM is able to be suspended
	req, err := http.NewRequest("PUT", config.ServerURL+path, body)
	if err != nil {
		return nil, err
	}

	//
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// WriteTCP
func pipeToConnection(conn net.Conn, in io.Reader, out io.Writer) error {

	// set up notification channels so that when a ping check fails, we can disconnect the active console
	pingService := make(chan interface{})

	// pipe data from the server to out, and from in to the server
	go func() {
		io.Copy(conn, in)
		conn.(*net.TCPConn).CloseWrite()
	}()

	//
	go func() {
		io.Copy(out, conn)
		close(pingService)
	}()

	//
	return monitorServer(pingService, 5*time.Second)
}

// monitor the server for disconnects
func monitorServer(done chan interface{}, wait time.Duration) error {

	ping := make(chan interface{}, 1)

	//
	for {

		// ping the server
		go func() {
			if ok, _ := Ping(); ok {
				ping <- true
			} else {
				close(ping)
			}
		}()

		select {

		//
		case _, ok := <-ping:

			//
			if !ok {
				return DisconnectedFromServer
			}

			//
			time.Sleep(time.Second)

			//
		case <-time.After(wait):
			return DisconnectedFromServer

		//
		case <-done:
			return nil
		}
	}
}
