//
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/nanobox-io/nanobox/config"
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
func connect(path string) (net.Conn, []byte, error) {

	//
	b := make([]byte, 1)

	// if we can't connect to the server, lets bail out early
	conn, err := net.Dial("tcp4", config.ServerURI)
	if err != nil {
		return conn, b, err
	}
	// we dont defer a conn.Close() here because we're returning the conn and
	// want it to remain open

	// make an http request
	if _, err := fmt.Fprint(conn, path); err != nil {
		return conn, b, err
	}

	// wait trying to read from the connection until a single read happens (blocking)
	if _, err := conn.Read(b); err != nil {
		return conn, b, err
	}

	return conn, b, nil
}

// pipeToConnection
func pipeToConnection(conn net.Conn, in io.Reader, out io.Writer) error {

	// set up notification channels so that when a ping check fails, we can disconnect the active console
	pong := make(chan interface{}, 1)
	done := make(chan interface{})

	// pipe data from the conn (server) to in
	go func() {
		io.Copy(conn, in)
		conn.(*net.TCPConn).CloseWrite()
	}()

	// pipe data from out to the conn (server)
	go func() {
		io.Copy(out, conn)
		close(done)
	}()

	// monitor the server connection
	for {

		// ping the server
		go func() {
			if ok, err := ping(); !ok {
				fmt.Printf("Ping failed and the server returned: %s", err.Error())
				close(pong)
				return
			}

			pong <- true
		}()

		// wait for the result of pinging the server
		select {

		//
		case _, ok := <-pong:

			// if the ping fails pong is closed
			if !ok {
				return fmt.Errorf("Connection lost... Server returned no response to ping.")
			}

			// loop, pinging at 1 second intervals
			<-time.After(time.Second * 1)

		// ping failed after 5 seconds
		case <-time.After(5 * time.Second):
			return fmt.Errorf("Connection timeout... Server returned no response to ping after 5 seconds.")

			//
		case <-done:
			return nil
		}
	}
}

// ping issues a ping to nanobox server
func ping() (bool, error) {

	// a new client is used to allow for shortening the request timeout
	client := http.Client{Timeout: time.Duration(2 * time.Second)}

	//
	res, err := client.Get(config.ServerURL + "/ping")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	//
	return res.StatusCode/100 == 2, nil
}
