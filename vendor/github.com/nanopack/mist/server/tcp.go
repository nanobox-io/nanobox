package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/jcelliott/lumber"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/core"
)

// init adds "tcp" as an available mist server type
func init() {
	Register("tcp", StartTCP)
}

// StartTCP starts a tcp server listening on the specified address (default 127.0.0.1:1445)
// and then continually reads from the server handling any incoming connections
func StartTCP(uri string, errChan chan<- error) {

	// start a TCP listener
	ln, err := net.Listen("tcp", uri)
	if err != nil {
		errChan <- fmt.Errorf("Failed to start tcp listener - %v", err.Error())
		return
	}

	lumber.Info("TCP server listening at '%s'...", uri)

	// start continually listening for any incoming tcp connections (non-blocking)
	go func() {
		for {

			// accept connections
			conn, err := ln.Accept()
			if err != nil {
				errChan <- fmt.Errorf("Failed to accept TCP connection %v", err.Error())
				return
			}

			// handle each connection individually (non-blocking)
			go handleConnection(conn, errChan)
		}
	}()
}

// handleConnection takes an incoming connection from a mist client (or other client)
// and sets up a new subscription for that connection, and a 'publish Handler'
// that is used to publish messages to the data channel of the subscription
func handleConnection(conn net.Conn, errChan chan<- error) {

	// close the connection when we're done here
	defer conn.Close()

	// create a new client for each connection
	proxy := mist.NewProxy()
	defer proxy.Close()

	// add basic TCP command handlers for this connection
	handlers = GenerateHandlers()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	// publish mist messages (pong, etc.. and messages if subscriber attatched)
	// to connected tcp client (non-blocking)
	go func() {
		for msg := range proxy.Pipe {
			lumber.Info("Got message - %#v", msg)
			// if the message fails to encode its probably a syntax issue and needs to
			// break the loop here because it will never be able to encode it; this will
			// disconnect the client.
			if err := encoder.Encode(msg); err != nil {
				errChan <- fmt.Errorf("Failed to pubilsh proxy.Pipe contents to TCP clients - %v", err)
				break
			}
		}
	}()

	// connection loop (blocking); continually read off the connection. Once something
	// is read, check to see if it's a message the client understands to be one of
	// its commands. If so attempt to execute the command.
	for {
		msg := mist.Message{}

		// if the message fails to decode its probably a syntax issue and needs to
		// break the loop here because it will never be able to decode it; this will
		// disconnect the client.
		if err := decoder.Decode(&msg); err != nil {
			switch err {
			case io.EOF:
				lumber.Debug("Client disconnected")
			case io.ErrUnexpectedEOF:
				lumber.Debug("Client disconnected unexpedtedly")
			default:
				errChan <- fmt.Errorf("Failed to decode message from TCP connection - %v", err)
			}
			return
		}

		// if an authenticator was passed, check for a token on connect to see if
		// auth commands are allowed
		if auth.DefaultAuth != nil && !proxy.Authenticated {

			// if the next input does not match the token then
			if msg.Data != authtoken {
				lumber.Debug("Client data doesn't match configured auth token")
				// break // allow connection w/o admin commands
				return // disconnect client
			}

			// todo: is this still used?
			// add auth commands ("admin" mode)
			for k, v := range auth.GenerateHandlers() {
				handlers[k] = v
			}

			// establish that the connection has already authenticated
			proxy.Authenticated = true
		}

		// look for the command
		handler, found := handlers[msg.Command]

		// if the command isn't found, return an error and wait for the next command
		if !found {
			lumber.Trace("Command '%v' not found", msg.Command)
			encoder.Encode(&mist.Message{Command: msg.Command, Tags: msg.Tags, Data: msg.Data, Error: "Unknown Command"})
			continue
		}

		// attempt to run the command; if the command fails return the error and wait
		// for the next command
		lumber.Trace("TCP Running '%v'...", msg.Command)
		if err := handler(proxy, msg); err != nil {
			lumber.Debug("TCP Failed to run '%v' - %v", msg.Command, err)
			encoder.Encode(&mist.Message{Command: msg.Command, Error: err.Error()})
			continue
		}
	}
}
