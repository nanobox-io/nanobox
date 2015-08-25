// -*- mode: go; tab-width: 2; indent-tabs-mode: 1; st-rulers: [70] -*-
// vim: ts=4 sw=4 ft=lua noet
//--------------------------------------------------------------------
// @author Daniel Barney <daniel@nanobox.io>
// @copyright 2015, Pagoda Box Inc.
// @doc
//
// @end
// Created :   12 August 2015 by Daniel Barney <daniel@nanobox.io>
//--------------------------------------------------------------------
package mist

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type (
	command struct {
		Command string `json:"command"`
	}
	tagList struct {
		Tags []string `json:"tags"`
	}
	list struct {
		Subscriptions [][]string `json:"subscriptions"`
		Command       string     `json:"command"`
		Success       bool       `json:"success"`
	}
)

//
func GenerateWebsocketUpgrade(mist *Mist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// we don't want this to be buffered
		client := NewLocalClient(mist, 0)

		write := make(chan string)
		done := make(chan bool)
		defer func() {
			client.Close()
			close(done)
		}()

		// the gorilla websocket package must have all writes come from the
		// same process.
		go func() {
			for {
				select {
				case event := <-client.Messages():
					if msg, err := json.Marshal(event); err == nil {
						conn.WriteMessage(websocket.TextMessage, msg)
					}
				case msg := <-write:
					conn.WriteMessage(websocket.TextMessage, []byte(msg))
				case <-done:
					close(write)
					return
				}
			}
		}()

		for {
			messageType, frame, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if messageType != websocket.TextMessage {
				write <- "{\"success\":false,\"error\":\"I don't understand binary messages\"}"
				continue
			}

			cmd := command{}
			if err := json.Unmarshal(frame, &cmd); err != nil {
				write <- "{\"success\":false,\"error\":\"Invalid json\"}"
				continue
			}

			switch cmd.Command {
			case "subscribe":
				tags := tagList{}
				// error would already be caught by unmarshalling the command
				json.Unmarshal(frame, &tags)
				client.Subscribe(tags.Tags)
				write <- "{\"success\":true,\"command\":\"subscribe\"}"
			case "unsubscribe":
				tags := tagList{}
				// error would already be caught by unmarshalling the command
				json.Unmarshal(frame, &tags)
				client.Unsubscribe(tags.Tags)
				write <- "{\"success\":true,\"command\":\"unsubscribe\"}"
			case "list":
				list := list{}
				list.Subscriptions, err = client.List()
				if err != nil {
					// do we need to do something with this error?
					return
				}
				list.Command = "list"
				list.Success = true
				bytes, err := json.Marshal(list)
				if err != nil {
					// Do I need to do something more here?
					return
				}
				write <- string(bytes)
			case "ping":
				write <- "{\"success\":true,\"command\":\"ping\"}"
			default:
				write <- "{\"success\":false,\"error\":\"unknown command\"}"
			}
		}
	}
}
