// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"strings"

	"github.com/pagodabox/nanobox-golang-stylish"
)

// start starts a tcp server listening on the specified port (default 1445), it
// then continually reads from the server handling any incoming connections
func (m *Mist) start() {
	m.log.Info(stylish.Bullet("Starting mist server..."))

	//
	go func() {

		//
		l, err := net.Listen("tcp", ":"+m.port)
		if err != nil {
			m.log.Error("%+v\n", err)
		}

		defer l.Close()

		m.log.Info(stylish.Bullet("Mist listening on port " + m.port))

		// Continually listen for any incoming connections.
		for {
			conn, err := l.Accept()
			if err != nil {
				m.log.Error("%+v\n", err)
			}

			// handle each connection individually (non-blocking)
			go m.handleConnection(conn)
		}
	}()
}

// handleConnection takes an incoming connection from a mist client (or other client)
// and sets up a new subscription for that connection, and a 'publish handler'
// that is used to publish messages to the data channel of the subscription
func (m *Mist) handleConnection(conn net.Conn) {
	m.log.Debug("[MIST :: SERVER] New connection detected: %+v\n", conn)

	// create a new subscription
	sub := Subscription{
		Sub: make(chan Message),
	}

	// make a done channel
	done := make(chan bool)

	// create a 'publish handler'
	go func() {
		for {

			// when a message is recieved on the subscriptions channel, append the length
			// of the message into the first 4 bytes so clients can know how big of a
			// message they should expect, and then write the message to the connection
			select {
			case msg := <-sub.Sub:

				b, err := json.Marshal(msg)
				if err != nil {
					m.log.Error("[MIST :: SERVER] Failed to marshal message: %v\n", err)
				}

				//
				bsize := make([]byte, 4)
				binary.LittleEndian.PutUint32(bsize, uint32(len(b)))

				if _, err := conn.Write(append(bsize, b...)); err != nil {
					break
				}

			// once the server is done sending messages issue a 'done'
			case <-done:
				break

			}
		}
	}()

	//
	r := bufio.NewReader(conn)

	//
	for {

		// read messages coming across the tcp channel
		l, err := r.ReadString('\n')
		if err != nil {

			// if communication stops, close the connection, unsubscribe, and issue 'done'
			if err == io.EOF {
				conn.Close()
				m.Unsubscribe(sub)
				done <- true

				// the channel is not closed here, because this is left up to the client
				// close(sub.Sub)
				break

				// some unexpected error happened
			} else {
				m.log.Error("[MIST :: SERVER] Error reading stream: %+v\n", err.Error())
			}
		}

		split := strings.Split(strings.TrimSpace(l), " ")
		cmd := split[0]

		//
		switch cmd {
		case "subscribe":
			sub.Tags = strings.Split(split[1], ",")
			m.Subscribe(sub)
		case "unsubscribe":
			sub.Tags = strings.Split(split[1], ",")
			m.Unsubscribe(sub)
		case "subscriptions":
			m.List()
		default:
			m.log.Error("[MIST :: SERVER] Unknown command: %+v\n", cmd)
		}
	}

	return
}
