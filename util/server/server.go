//
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
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
	return http.Post(config.ServerURL+path, contentType, body)
}

// Put
func Put(path string, body io.Reader) (*http.Response, error) {

	// check to see if the VM is able to be suspended
	req, err := http.NewRequest("PUT", config.ServerURL+path, body)
	if err != nil {
		return nil, err
	}

	//
	return http.DefaultClient.Do(req)
}

// connect
func connect(params string) (net.Conn, error) {

	// if we can't connect to the server, lets bail out early
	conn, err := net.Dial("tcp4", config.ServerURI)
	if err != nil {
		return conn, err
	}
	defer conn.Close()

	// make a http request
	if _, err := fmt.Fprintf(conn, "POST /exec?pid=%d&%v HTTP/1.1\r\n\r\n", os.Getpid(), params); err != nil {
		return conn, err
	}

	//
	p := make([]byte, 512)

	// wait trying to read from the connection until a single read happens (blocking)
	for {
		n, err := conn.Read(p)
		if err != nil {
			return conn, err
		}

		// once a single read happens write the string and break out of the loop
		if n > 0 {
			os.Stderr.WriteString(string(p))
			break
		}
	}

	return conn, nil
}

// pipeToConnection
func pipeToConnection(conn net.Conn, in io.Reader, out io.Writer) error {

	// set up notification channels so that when a ping check fails, we can disconnect the active console
	pingService := make(chan interface{})
	ping := make(chan interface{}, 1)

	// pipe data from the server to out
	go func() {
		io.Copy(conn, in)
		conn.(*net.TCPConn).CloseWrite()
	}()

	// pipe data from in to the server
	go func() {
		io.Copy(out, conn)
		close(pingService)
	}()

	// monitor the server connection
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
		case <-time.After(5 * time.Second):
			return DisconnectedFromServer

		//
		case <-pingService:
			return nil
		}
	}
}
