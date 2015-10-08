// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package mist

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/nanobox-io/golang-mist"
	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
type (

	// Model
	Model struct {
		Action   string `json:"action"`
		Document struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"document"`
		Name string `json:"model"`
	}

	// Log
	Log struct {
		Content  string `json:"content"`
		Priority int    `json:"priority"`
		Time     string `json:"time"`
		Type     string `json:"type"`
	}
)

//
var (

	// subscriptions is a list of tags that have been used to subscribe with either
	// Listen or Stream; when creating a new Listner/Streamer if the tags have
	// already been used, it stops double subscription
	subscriptions = make(map[string]struct{})

	// a map of each type of 'process' that we encounter to then be used when
	// assigning a unique color to that 'process'
	logProcesses = make(map[string]string)

	// an array of the colors used to colorize the logs
	logColors = [11]string{
		// "red",
		"green",
		"yellow",
		"blue",
		"magenta",
		"cyan",
		// "light_red", // this is reserved for a failover output
		"light_green",
		"light_yellow",
		"light_blue",
		"light_magenta",
		"light_cyan",
		"white",
	}
)

// Listen connects a to mist, subscribes tags, and listens for 'model' updates
func Listen(tags []string, handle func(string) bool) error {

	// if this subscription already exists, exit; this prevents double subscriptions
	if _, ok := subscriptions[strings.Join(tags, "")]; ok {
		return nil
	}

	// connect
	client, err := connect()
	if err != nil {
		fmt.Printf(stylish.ErrBullet("Failed to create client - %s", err.Error()))
	}
	defer client.Close()

	// subscribe
	if err := subscribe(client, tags); err != nil {
		fmt.Printf(stylish.ErrBullet("Failed to subscribe to %v - %s", tags, err.Error()))
	}

	// add tags to list of subscriptions
	subscriptions[strings.Join(tags, "")] = struct{}{}

	//
	model := Model{}
	for msg := range client.Messages() {

		// unmarshal the incoming Message
		if err := json.Unmarshal([]byte(msg.Data), &model); err != nil {
			return err
		}

		// handle the status; when the handler returns false, it's time to break the
		// stream
		if listen := handle(model.Document.Status); !listen {
			return nil
		}
	}

	return nil
}

// Stream connects to mist, subscribes tags, and logs Messages
func Stream(tags []string, handle func(Log)) {

	// add log level to tags
	tags = append(tags, config.LogLevel)

	// if this subscription already exists, exit; this prevents double subscriptions
	if _, ok := subscriptions[strings.Join(tags, "")]; ok {
		return
	}

	// connect
	client, err := connect()
	if err != nil {
		fmt.Printf(stylish.ErrBullet("Failed to create client - %s", err.Error()))
	}
	defer client.Close()

	// subscribe
	if err := subscribe(client, tags); err != nil {
		fmt.Printf(stylish.ErrBullet("Failed to subscribe to %v - %s", tags, err.Error()))
	}

	// add tags to list of subscriptions
	subscriptions[strings.Join(tags, "")] = struct{}{}

	//
	for msg := range client.Messages() {

		//
		log := Log{}

		// unmarshal the incoming Message
		if err := json.Unmarshal([]byte(msg.Data), &log); err != nil {
			config.Fatal("[util/server/mist] failed - ", err.Error())
		}

		//
		handle(log)
	}
}

// connect connects 'mist' to the server running on the guest machine
func connect() (mist.Client, error) {
	return mist.NewRemoteClient(config.MistURI)
}

// subscribe
func subscribe(client mist.Client, tags []string) error {
	return client.Subscribe(tags)
}

// ProcessLog takes a Logvac or Stormpack log and breaks it apart into pieces that
// are then reconstructed in a 'digestible' way, colorized, and output to the
// terminal
func ProcessLog(log Log) {

	// t := time.Now(log.Time).Format(time.RFC822)
	// t, err := time.Parse("01/02 03:04:05PM '06 -0700", log.Time)
	// if err != nil {
	// 	fmt.Println("TIME BONK!", err)
	// }

	//
	subMatch := regexp.MustCompile(`^(\w+)\.(\S+)\s+(.*)$`).FindStringSubmatch(log.Content)

	// ensure a subMatch and ensure subMatch has a length of 4, since thats how many
	// matches we're expecting
	if subMatch != nil && len(subMatch) >= 4 {

		service := subMatch[1]
		process := subMatch[2]
		content := subMatch[3]

		//
		if _, ok := logProcesses[process]; !ok {
			logProcesses[process] = logColors[len(logProcesses)%len(logColors)]
		}

		// util.Printc("[%v]%v - %v.%v :: %v[reset]", logProcesses[process], log.Time, service, process, content)
		util.Printc("[%v]%v (%v) :: %v[reset]", logProcesses[process], service, process, content)

		// if we don't have a subMatch or its length is less than 4, just print w/e
		// is in the log
	} else {
		util.Printc("[light_red]%v - %v[reset]", log.Time, log.Content)
	}

}
