// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

// build a small applicationController so that we don't have to play with a large
// switch statement
type (
	handler struct {
		argCount int
		handle   func(*localSubscriber, []string) string
	}
)

var (
	commandMap = map[string]handler{
		"list":        {0, handleList},
		"subscribe":   {1, handleSubscribe},
		"unsubscribe": {1, handleUnubscribe},
		"publish":     {2, handlePublish},
		"ping":        {0, handlePing},
	}
)

//
func handlePing(client *localSubscriber, args []string) string {
	return "pong\n"
}

//
func handleSubscribe(client *localSubscriber, args []string) string {
	tags := strings.Split(args[0], ",")
	client.Subscribe(tags)
	return ""
}

//
func handleUnubscribe(client *localSubscriber, args []string) string {
	tags := strings.Split(args[0], ",")
	client.Unsubscribe(tags)
	return ""
}

//
func handleList(client *localSubscriber, args []string) string {
	list, err := client.List()
	if err != nil {
		return err.Error()
	}
	tmp := make([]string, len(list))

	for idx, subscription := range list {
		tmp[idx] = strings.Join(subscription, ",")
	}

	response := strings.Join(tmp, " ")
	return fmt.Sprintf("list %v\n", response)
}

//
func handlePublish(client *localSubscriber, args []string) string {
	tags := strings.Split(args[0], ",")
	client.Publish(tags, args[1])
	return ""
}

// start starts a tcp server listening on the specified address (default 127.0.0.1:1445),
// it then continually reads from the server handling any incoming connections
func (m *Mist) Listen(address string) (net.Listener, error) {
	if address == "" {
		address = "127.0.0.1:1445"
	}
	serverSocket, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	go func() {
		defer serverSocket.Close()
		// Continually listen for any incoming connections.
		for {
			conn, err := serverSocket.Accept()
			if err != nil {
				// what should we do with the error?
				return
			}

			// handle each connection individually (non-blocking)
			go m.handleConnection(conn)
		}
	}()
	return serverSocket, nil
}

// handleConnection takes an incoming connection from a mist client (or other client)
// and sets up a new subscription for that connection, and a 'publish handler'
// that is used to publish messages to the data channel of the subscription
func (m *Mist) handleConnection(conn net.Conn) {

	// create a new client to match with this connection

	client := NewLocalClient(m, 0)

	// make a done channel
	done := make(chan bool)

	// clean up everything that we have setup
	defer func() {
		conn.Close()
		client.Close()
		close(done)
	}()

	// create a 'publish handler'
	go func() {
		for {

			// when a message is recieved on the subscriptions channel write the message
			// to the connection
			select {
			case msg := <-client.Messages():

				if _, err := conn.Write([]byte(fmt.Sprintf("publish %v %v\n", strings.Join(msg.Tags, ","), msg.Data))); err != nil {
					break
				}

			// return if we are done
			case <-done:
				return
			}
		}
	}()

	//
	r := bufio.NewReader(conn)

	//
	for {

		// read messages coming across the tcp channel
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			// some unexpected error happened
			return
		}

		line = strings.TrimSuffix(line, "\n")

		// this is the general format of the commands that are accepted
		// ["cmd" ,"tag,tag2", "all the rest"]
		split := strings.SplitN(line, " ", 3)
		cmd := split[0]

		handler, found := commandMap[cmd]

		var response string
		args := split[1:]

		switch {
		case !found:
			response = fmt.Sprintf("error Unknown command '%v'\n", cmd)
		case handler.argCount != len(args):
			response = fmt.Sprintf("error Incorrect number of arguments for '%v'\n", cmd)
		default:
			response = handler.handle(client, args)
		}

		if response != "" {
			// Is it safe to send from 2 gorountines at the same time?
			if _, err := conn.Write([]byte(response)); err != nil {
				break
			}
		}
	}
}
