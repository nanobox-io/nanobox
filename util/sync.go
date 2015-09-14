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
		Action   string `json:"action"`
		Document string `json:"document"`
		Log      string `json:"log"`
		Model    string `json:"model"`
		Time     string `json:"time"`
	}
)

// run issues a sync to the running nanobox VM
func (s *Sync) Run(opts []string) {

	// connect 'mist' to the server running on the guest machine
	client, err := mist.NewRemoteClient(config.MistURI)
	if err != nil {
		LogFatal("[utils/sync] client.Connect() failed ", err)
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
		LogFatal("[utils/sync] api.DoRawRequest() failed ", err)
	}

	// handle
stream:
	for msg := range client.Messages() {

		//
		e := &entry{}

		// check for any error message that cause mist to disconnect, and release
		// nanobox
		// if msg.Tags[0] == "err" {
		// 	fmt.Printf(stylish.Warning("An unexpected error caused the sync stream to disconnect. Your sync should continue as normal, and you can see the log output from your dashboard."))
		// 	break stream
		// }

		//
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			LogFatal("[utils/sync] json.Unmarshal() failed", err)
		}

		// depending on what fields the data has, determines what needs to happen...
		switch {

		// if the message contains the log field, the log is printed. The message is
		// then checked to see if it contains a model field...
		// example entry: {Time: "time", Log: "content"}
		case e.Log != "":
			fmt.Printf(e.Log)

			// if the message contains the model field...
		case e.Model != "":

			// update the model status
			if err := json.Unmarshal([]byte(e.Document), s); err != nil {
				LogFatal("[utils/sync] json.Unmarshal() failed ", err)
			}

			// break the stream once we get a model update. If we ever have intermediary
			// statuses we can throw in a case that will handle this on a status-by-status
			// basis (errored, complete)
			break stream

		// report any unhandled entries, incase cases need to be added to handle them
		default:
			config.Console.Debug("Nanobox has encountered an Entry with neither a 'log' nor 'model' field and doesnt know what to do...")
			break stream
		}
	}
}
