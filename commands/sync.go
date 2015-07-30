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
	"strings"

	"github.com/pagodabox/golang-mist"
	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

type (

	// Sync
	nsync struct {
		kind    string
		path    string
		verbose bool

		ID     string `json:"id"`
		Status string `json:"status"`
	}

	// Entry
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

	// subscribe
	// create a 'mist' client to communicate with the mist server running on the
	// guest machine
	client := mist.Client{Host: config.Nanofile.IP, Port: "1445"}

	//
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

			// depending on the type of model, different things may happen...
			switch e.Model {

			// in the case of a sync model, listen for a complete to close the stream
			case strings.Title(s.kind), s.kind:

				if err := json.Unmarshal([]byte(e.Document), s); err != nil {
					ui.LogFatal("[commands sync] json.Unmarshal() failed ", err)
				}

				switch s.Status {
				// once the sync is 'complete' unsubscribe from mist
				case "complete":
					fmt.Printf(stylish.Bullet(fmt.Sprintf("%v complete... Navigate to %v.nano.dev to view your app.", strings.Title(s.kind), config.App)))
					break stream

				// if the sync is 'errored' unsubscribe from mist
				case "errored":
					fmt.Printf(stylish.Error(fmt.Sprintf("%v failed", strings.Title(s.kind)), fmt.Sprintf("Your %v failed to uh... %v", s.kind, s.kind)))
					break stream
				}

			// report any unhandled models, incase cases need to be added to handle them
			case "default":
				config.Console.Debug("Nanobox has encountered an unknown model (%v), and doesn't know what to do with it...", e.Model)
				break stream
			}

		// report any unhandled entries, incase cases need to be added to handle them
		default:
			config.Console.Debug("Nanobox has encountered an Entry with neither a 'log' nor 'model' field and doesnt know what to do...")
			break stream
		}
	}
}
