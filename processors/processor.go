package processors

import (
	"bytes"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
)

var (

	// used for console and tunnel
	container string
	key       string
	location  string
)

// getAppID ...
func getAppID(alias string) string {
	env, _ := models.FindEnvByID(config.EnvID())
	if alias == "" {
		alias = "default"
	}
	app, ok := env.Links[alias]
	if !ok {
		return alias
	}

	return app
}

// connect ...
func connect(req *http.Request) (net.Conn, *bytes.Buffer, error) {
	//
	b := make([]byte, 1)

	// if we can't connect to the server, lets bail out early
	conn, err := tls.Dial("tcp4", location, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return conn, bytes.NewBuffer(b), err
	}


	// we dont defer a conn.Close() here because we're returning the conn and
	// want it to remain open

	// make an http request

	if err := req.Write(conn); err != nil {
		return conn, bytes.NewBuffer(b), err
	}

	return conn, bytes.NewBuffer(b), nil
}
