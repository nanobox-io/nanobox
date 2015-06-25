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
	"os"

	mist "github.com/pagodabox/golang-mist"
	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

type (

	// DeployCommand satisfies the Command interface for listing a user's apps
	DeployCommand struct{}

	// Sync
	Sync struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
)

// Help prints detailed help text for the app list command
func (c *DeployCommand) Help() {
	ui.CPrint(`
Description:
  Deploys to your nanobox VM

Usage:
  nanobox deploy
  `)
}

// Run displays select information about all of a user's apps
func (c *DeployCommand) Run(opts []string) {

	// start the vm if it's not already running
	// resume := ResumeCommand{}
	// resume.Run(opts)

	// subscribe to mist
	client := mist.Client{}
	if _, err := client.Connect(config.Boxfile.IP, "1445"); err != nil {
		ui.LogFatal("[commands deploy] client.Connect() failed ", err)
	}

	defer client.Close()

	sub, err := client.Subscribe([]string{"sync", "deploy"})
	if err != nil {
		config.Console.Warn("Failed to subscribe to 'mist' updates... %v", err)
	}

	// issue a deploy
	path := fmt.Sprintf("http://%v:1757/deploys", config.Boxfile.IP)

	if err := api.DoRawRequest(nil, "POST", path, nil, nil); err != nil {
		ui.LogFatal("[commands deploy] api.DoRawRequest() failed ", err)
	}

	// listen for messages coming from mist
	for msg := range client.Data {

		data := make(map[string]string)

		//
		if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
			ui.LogFatal("[commands deploy] json.Unmarshal() failed ", err)
		}

		// depending on what fields the data has, determines what needs to happen
		switch {

		// if the message contains the log field, the log is printed...
		case data["log"] != "":
			entry := fmt.Sprintf("[%v] %v", data["log"], data["time"])

			fmt.Println(entry)

		// if the message contains the model field, handle individually
		case data["model"] != "":

			// depending on the type of model, different things may happen...
			switch data["model"] {

			// in the case of a sync model, listen for a complete to close the stream
			case "sync":
				sync := &Sync{}

				if err := json.Unmarshal([]byte(data["document"]), sync); err != nil {
					ui.LogFatal("[commands deploy] json.Unmarshal() failed ", err)
				}

				// once the sync is 'complete' unsubscribe from mist, and close the connection
				if sync.Status == "complete" {
					client.Unsubscribe(sub)
					client.Close()
				}

			// the only type of model expected for the time being is a sync, anything
			// else should fail because logic is probably missing to handle the new
			// model
			case "default":
				config.Console.Error("[commands deploy] Unhandled model '%v'", data["model"])
				os.Exit(1)
			}

		// if the message does not cotain either a 'log' or 'model' field the CLI
		// needs to fail, because it's probably missing some logic to handle a new
		// field
		default:
			config.Console.Error("[commands deploy] Unhandled data, missing 'log' or 'model': %v", data)
			os.Exit(1)
		}
	}

}
