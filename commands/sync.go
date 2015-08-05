// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"encoding/json"
	"fmt"

	"github.com/pagodabox/golang-mist"
	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

type (

	// nsync
	nsync struct {
		kind    string
		path    string
		verbose bool

		ID     string `json:"id"`
		Status string `json:"status"`
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
func (s *nsync) run(opts []string) {

	// start the vm if it's not already running
	// resume := ResumeCommand{}
	// resume.Run(opts)

	// create a 'mist' client to communicate with the mist server running on the
	// guest machine
	client := mist.Client{Host: config.Nanofile.IP, Port: "1445"}

	// connect the 'mist' client to the 'mist' server
	if err := client.Connect(); err != nil {
		ui.LogFatal("[commands sync] client.Connect() failed ", err)
	}
	defer client.Close()

	utils.Printv(stylish.Bullet("Subscribing to mist..."), s.verbose)

	// subscribe to 'sync' updates
	utils.Printv(stylish.SubBullet("- Subscribing to app logs"), s.verbose)
	jobSub, err := client.Subscribe([]string{"job", s.kind})
	if err != nil {
		fmt.Printf(stylish.Warning("Nanobox failed to subscribe to app logs. Your sync will continue as normal, and log output is available on your dashboard."))
	}
	defer client.Unsubscribe(jobSub)

	logLevel := "info"
	if s.verbose {
		logLevel = "debug"
	}

	// if the verbose flag is included, also subscribe to the 'debug' sync logs
	utils.Printv(stylish.SubBullet("- Subscribing to debug logs"), s.verbose)
	logSub, err := client.Subscribe([]string{"log", "deploy", logLevel})
	if err != nil {
		fmt.Printf(stylish.Warning("Nanobox failed to subscribe to debug logs. Your sync will continue as normal, and log output is available on your dashboard."))
	}
	defer client.Unsubscribe(logSub)

	utils.Printv(stylish.Success(), s.verbose)

	//
	// issue a sync
	if err := api.DoRawRequest(nil, "POST", s.path, nil, nil); err != nil {
		ui.LogFatal("[commands sync] api.DoRawRequest() failed ", err)
	}

	// handle
stream:
	for msg := range client.Data {

		//
		e := &entry{}

		// check for any error message that cause mist to disconnect, and release
		// nanobox
		if msg.Tags[0] == "err" {
			fmt.Printf(stylish.Warning("An unexpected error caused the sync stream to disconnect. Your sync should continue as normal, and you can see the log output from your dashboard."))
			break stream
		}

		//
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			ui.LogFatal("[commands sync] json.Unmarshal() failed", err)
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
				ui.LogFatal("[commands sync] json.Unmarshal() failed ", err)
			}

			// break the stream once we get a model update. If we ever have intermediary
			// status's we can throw in a case that will handle this on a status-by-status
			// basis
			break stream

		// report any unhandled entries, incase cases need to be added to handle them
		default:
			config.Console.Debug("Nanobox has encountered an Entry with neither a 'log' nor 'model' field and doesnt know what to do...")
			break stream
		}
	}
}
