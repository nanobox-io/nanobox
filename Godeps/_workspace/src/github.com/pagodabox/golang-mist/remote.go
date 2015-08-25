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
	"net"
	"strings"
)

//
type (

	// A remoteSubscriber represents a connection to the mist server
	remoteSubscriber struct {
		conn net.Conn        // the connection the mist server
		done chan bool       // the channel to indicate that the connection is closed
		pong chan bool       // the channel for ping responses
		list chan [][]string // the channel for subscription listing
		data chan Message    // the channel that mist server 'publishes' updates to
	}
)

// Connect attempts to connect to a running mist server at the clients specified
// host and port.
func NewRemoteClient(address string) (Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	client := remoteSubscriber{
		done: make(chan bool),
		pong: make(chan bool),
		list: make(chan [][]string),
	}
	client.conn = conn

	// create a channel on which to publish messages received from mist server
	client.data = make(chan Message)

	// continually read from conn, forwarding the data onto the clients data channel
	go func() {
		defer close(client.data)

		r := bufio.NewReader(client.conn)
		for {
			var listChan chan [][]string
			var pongChan chan bool
			var dataChan chan Message

			line, err := r.ReadString('\n')
			if err != nil {
				// do we need to log the error?
				return
			}
			line = strings.TrimSuffix(line, "\n")

			// create a new message
			var msg Message
			var list [][]string

			split := strings.SplitN(line, " ", 2)

			switch split[0] {
			case "publish":
				split := strings.SplitN(split[1], " ", 2)
				msg = Message{
					Tags: strings.Split(split[0], ","),
					Data: split[1],
				}
				dataChan = client.data
			case "pong":
				pongChan = client.pong
			case "list":
				split := strings.Split(split[1], " ")
				list = make([][]string, len(split))
				for idx, subscription := range split {
					list[idx] = strings.Split(subscription, ",")
				}
				listChan = client.list
			case "error":
				// need to report the error somehow
				// close the connection as something is seriously wrong
				client.Close()
				return
			}

			// send the message on the client channel, or close if this connection is done
			select {
			case listChan <- list:
			case pongChan <- true:
			case dataChan <- msg:
			case <-client.done:
				return
			}
		}
	}()

	return &client, nil
}

// List requests a list of current mist subscriptions from the server
func (client *remoteSubscriber) List() ([][]string, error) {
	if _, err := client.conn.Write([]byte("list\n")); err != nil {
		return nil, err
	}
	return <-client.list, nil
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (client *remoteSubscriber) Subscribe(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	_, err := client.conn.Write([]byte("subscribe " + strings.Join(tags, ",") + "\n"))

	return err
}

// Unsubscribe takes the specified tags and tells the server to unsubscribe from
// updates on those tags, returning an error or nil
func (client *remoteSubscriber) Unsubscribe(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	_, err := client.conn.Write([]byte("unsubscribe " + strings.Join(tags, ",") + "\n"))

	return err
}

// Publish sends a message to the mist server to be published to all subscribed clients
func (client *remoteSubscriber) Publish(tags []string, data string) error {
	if len(tags) == 0 {
		return nil
	}
	_, err := client.conn.Write([]byte(fmt.Sprintf("publish %v %v\n", strings.Join(tags, ","), data)))

	return err
}

// Ping pong the server
func (client *remoteSubscriber) Ping() error {
	if _, err := client.conn.Write([]byte("ping\n")); err != nil {
		return err
	}
	// wait for the pong to come back
	<-client.pong
	return nil
}

//
func (client *remoteSubscriber) Messages() <-chan Message {
	return client.data
}

// Close closes the client data channel and the connection to the server
func (client *remoteSubscriber) Close() error {
	// we need to do it in this order in case the goroutine is stuck waiting for
	// more data from the socket
	err := client.conn.Close()
	close(client.done)
	return err
}
