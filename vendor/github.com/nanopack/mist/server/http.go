package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/jcelliott/lumber"
)

var (

	// Router ...
	Router = pat.New()
)

// init adds http/https as available mist server types
func init() {
	Register("http", StartHTTP)
	Register("https", StartHTTPS)
}

// StartHTTP starts a mist server listening over HTTP
func StartHTTP(uri string, errChan chan<- error) {
	if err := newHTTP(uri); err != nil {
		errChan <- fmt.Errorf("Unable to start mist http listener - %v", err)
	}
}

// StartHTTPS starts a mist server listening over HTTPS
func StartHTTPS(uri string, errChan chan<- error) {
	errChan <- ErrNotImplemented
}

//
func newHTTP(address string) error {
	lumber.Info("HTTP server listening at '%s'...\n", address)

	// blocking...
	return http.ListenAndServe(address, routes())
}

// routes registers all api routes with the router
func routes() *pat.Router {

	//
	Router.Get("/ping", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("pong\n"))
	})
	// Router.Get("/list", handleRequest(list))
	// Router.Get("/subscribe", handleRequest(subscribe))
	// Router.Get("/unsubscribe", handleRequest(unsubscribe))

	return Router
}

// handleRequest is a wrapper for the actual route handler, simply to provide some
// debug output
func handleRequest(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		fn(rw, req)

		// must be after fn if ever going to get rw.status (logging still more meaningful)
		// config.Log.Trace(`%v - [%v] %v %v %v(%s) - "User-Agent: %s", "X-Nanobox-Token: %s"`,
		// 	req.RemoteAddr, req.Proto, req.Method, req.RequestURI,
		// 	rw.Header().Get("status"), req.Header.Get("Content-Length"),
		// 	req.Header.Get("User-Agent"), req.Header.Get("X-Nanobox-Token"))
	}
}
