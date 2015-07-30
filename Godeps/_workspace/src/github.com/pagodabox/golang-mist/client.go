// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/pagodabox/nanobox-golang-stylish"
)

//
type (

	// A Client has a connection to the mist server, a data channel from which it
	// receives messages from the server, ad a host and port to use when connecting
	// to the server
	Client struct {
		conn net.Conn     // the connection the mist server
		Data chan Message // the channel that mist server 'publishes' updates to
		Host string       // the connection host for where mist server is running
		Port string       // the connection port for where mist server is running
	}
)

// Connect attempts to connect to a running mist server at the clients specified
// host and port.
func (c *Client) Connect() error {
	fmt.Printf(stylish.Bullet("Attempting to connect to mist..."))

	// number of seconds/attempts to try when failing to conenct to mist server
	maxRetries := 60

	// attempt to connect to the host:port
	for i := 0; i < maxRetries; i++ {
		if conn, err := net.Dial("tcp", c.Host+":"+c.Port); err != nil {

			// max number of attempted retrys failed...
			if i >= maxRetries {
				fmt.Printf(stylish.Error("mist connection failed", "The attempted connection to mist failed. This shouldn't effect any running processes, however no output should be expected"))
				return err
			}
			fmt.Printf("\r   Connection failed! Retrying (%v/%v attempts)...", i, maxRetries)

			// upon successful connection, set the clients connection (conn) to the tcp
			// connection that was established with the server
		} else {
			fmt.Printf(stylish.SubBullet("- Connection established"))
			fmt.Printf(stylish.Success())
			c.conn = conn
			break
		}

		//
		time.Sleep(1 * time.Second)
	}

	// create a channel on which to publish messages received from mist server
	c.Data = make(chan Message)

	// continually read from conn, forwarding the data onto the clients data channel
	go func() {
		for {

			// read the first 4 bytes of the message so we know how long the message
			// is expected to be
			bsize := make([]byte, 4)
			if _, err := io.ReadFull(c.conn, bsize); err != nil {

				// for now, the only err I can see causing problems here is a closed
				// connection, so for now we'll just break until we need to handle
				// different types
				break
			}

			// create a buffer that is the length of the expected message
			n := binary.LittleEndian.Uint32(bsize)

			// read the length of the message up to the expected bytes
			b := make([]byte, n)
			if _, err := io.ReadFull(c.conn, b); err != nil {

				// for now, the only err I can see causing problems here is a closed
				// connection, so for now we'll just break until we need to handle
				// different types
				break
			}

			// create a new message
			msg := Message{}

			// unmarshal the raw message into a mist message
			if err := json.Unmarshal(b, &msg); err != nil {
				c.Data <- Message{Tags: []string{"err"}, Data: err.Error()}
			}

			// send the message on the client data channel to be handled from the clients
			// user
			c.Data <- msg
		}
	}()

	return nil
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (c *Client) Subscribe(tags []string) ([]string, error) {
	if _, err := c.conn.Write([]byte("subscribe " + strings.Join(tags, ",") + "\n")); err != nil {
		return nil, err
	}

	return tags, nil
}

// Unsubscribe takes the specified tags and tells the server to unsubscribe from
// updates on those tags, returning an error or nil
func (c *Client) Unsubscribe(tags []string) error {
	if _, err := c.conn.Write([]byte("unsubscribe " + strings.Join(tags, ",") + "\n")); err != nil {
		return err
	}

	return nil
}

// Subscriptions requests a list of current mist subscriptions from the server
func (c *Client) Subscriptions() error {
	if _, err := c.conn.Write([]byte("subscriptions\n")); err != nil {
		return err
	}

	return nil
}

// Close closes the client data channel and the connection to the server
func (c *Client) Close() error {
	close(c.Data)
	return c.conn.Close()
}
