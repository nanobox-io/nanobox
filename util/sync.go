// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pagodabox/golang-mist"
	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-golang-stylish"
)

type (

	// nsync
	Sync struct {
		ID      string `json:"id"`
		Model   string
		Path    string
		Status  string `json:"status"`
		Verbose bool
	}

	// entry
	entry struct {

		// 'model' message fields
		Action   string `json:"action"`
		Document string `json:"document"`
		Model    string `json:"model"`

		// 'log' message fields
		Content  string `json:"content"`
		Priority int    `json:"priority"`
		Time     string `json:"time"`
		Type     string `json:"type"`
	}
)

// run issues a sync to the running nanobox VM
func (s *Sync) Run(opts []string) {

	// connect 'mist' to the server running on the guest machine
	client, err := mist.NewRemoteClient(config.MistURI)
	if err != nil {
		config.Fatal("[utils/sync] client.Connect() failed ", err)
	}
	defer client.Close()

	Printv(stylish.Bullet("Subscribing to mist..."), s.Verbose)

	// subscribe to job updates
	jobTags := []string{"job", s.Model}
	if err := client.Subscribe(jobTags); err != nil {
		fmt.Printf(stylish.Warning("Nanobox failed to subscribe to app logs. Your sync will continue as normal, and log output is available on your dashboard."))
	}
	defer client.Unsubscribe(jobTags)

	Printv(stylish.SubBullet("- Subscribed to app logs"), s.Verbose)

	logLevel := "info"
	if s.Verbose {
		logLevel = "debug"
	}

	// if the verbose flag is included, also subscribe to the 'debug' logs
	logTags := []string{"log", "deploy", logLevel}
	if err := client.Subscribe(logTags); err != nil {
		fmt.Printf(stylish.Warning("Nanobox failed to subscribe to debug logs. Your sync will continue as normal, and log output is available on your dashboard."))
	}
	defer client.Unsubscribe(logTags)

	Printv(stylish.SubBullet("- Subscribed to debug logs"), s.Verbose)

	//
	// issue a sync
	if err := api.DoRawRequest(nil, "POST", s.Path, nil, nil); err != nil {
		config.Fatal("[utils/sync] api.DoRawRequest() failed", err)
	}

	// handle
stream:
	for msg := range client.Messages() {

		// unmarshal the incoming message; a message will be of two possible 'types'
		// 1. A log entry that has type, time, priority, and content
		// 2. A 'model' update with model, action, and document
		// once parsed, the following case statemt will determine what time of message
		// was received and what to do about it.
		e := &entry{}
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			config.Fatal("[utils/sync] json.Unmarshal() failed", err)
		}

		// depending on what fields the data has, determines what needs to happen...
		switch {

		// if the message contains the log field, the log is printed. The message is
		// then checked to see if it contains a model field...
		// example entry: {Time: "time", Log: "content"}
		case e.Content != "":
			fmt.Printf(e.Content)

			// if the message contains the model field...
		case e.Model != "":

			// update the model status
			if err := json.Unmarshal([]byte(e.Document), s); err != nil {
				config.Fatal("[utils/sync] json.Unmarshal() failed ", err)
			}

			switch s.Status {

			// for each successful deploy create/update the .nanobox/apps/<app>/.deployed
			// file
			case "complete":
				if _, err := os.Create(config.AppDir + "/.deployed"); err != nil {
					config.Fatal("[utils/sync] os.Create() failed ", err)
				}

			// case "unavailable":

			// if a deploy ever errors, remove the deployed file
			case "errored":
				if err := os.Remove(config.AppDir + "/.deployed"); err != nil {
					config.Fatal("[utils/sync] os.Remove() failed ", err)
				}
			}

			// break the stream once we get a model update. If we ever have intermediary
			// statuses we can throw in a case that will handle this on a status-by-status
			// basis (current statuses: complete, unavailable, errored)
			break stream

		// report any unhandled entries, incase cases need to be added to handle them
		default:
			fmt.Printf(stylish.ErrBullet("Malformed Entry..."))
			// break stream
		}
	}
}

// AppDeployed
func AppDeployed() bool {
	if _, err := os.Stat(config.AppDir + "/.deployed"); err != nil {
		return false
	}

	return true
}
