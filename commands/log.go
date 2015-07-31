// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/pagodabox/golang-mist"
	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	// "github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

type (

	// LogCommand satisfies the Command interface for obtaining an app's historical
	// and streaming logs
	LogCommand struct{}

	// Log represents the structure of a log returned from Logvac or Stormpack
	Log struct {
		Log  string `json:"log"`
		Time int    `json:"time"`
	}
)

// a map of each type of 'process' that we encounter to then be used when assigning
// a unique color to that 'process'
var logProcesses = make(map[string]string)

// an array of the colors used to colorize the logs
var logColors = [11]string{
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

// Help
func (c *LogCommand) Help() {
	ui.CPrint(`
Description:
  Provides the last 100 lines of historical log output by default.

Usage:
  pagoda log [-c] [-l]

Options:
  -c, --count
    The number of lines of the historical log you wish to view.

  -s, --stream
    Stream logs live

  -l --level
    Filters logs by one of the following levels:
			debug > info > warn > error > fatal

			note: that each level will display logs from all levels below it
  `)
}

// Run
func (c *LogCommand) Run(opts []string) {

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fCount int
	flags.IntVar(&fCount, "c", 100, "")
	flags.IntVar(&fCount, "count", 100, "")

	var fStream bool
	flags.BoolVar(&fStream, "s", false, "")
	flags.BoolVar(&fStream, "stream", false, "")

	var fLevel string
	flags.StringVar(&fLevel, "l", "info", "")
	flags.StringVar(&fLevel, "level", "info", "")

	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.app_log] flags.Parse()", err)
	}

	// if stream is true, we connect to the live logs
	if fStream {
		fmt.Printf(stylish.Bullet("Connecting to live stream..."))

		// create a 'mist' client to communicate with the mist server running on the
		// guest machine
		client := mist.Client{Host: config.Nanofile.IP, Port: "1445"}

		// connect the 'mist' client to the 'mist' server
		if err := client.Connect(); err != nil {
			ui.LogFatal("[commands sync] client.Connect() failed ", err)
		}
		defer client.Close()

		// subscribe to 'sync' updates
		logSub, err := client.Subscribe([]string{"log", fLevel})
		if err != nil {
			fmt.Printf(stylish.Warning("Nanobox failed to subscribe to app logs. Your sync will continue as normal, and log output is available on your dashboard."))
		}
		defer client.Unsubscribe(logSub)

		//
		fmt.Printf(stylish.Bullet("Connecting to live stream..."))

	stream:
		for msg := range client.Data {

			//
			log := Log{}

			// check for any error message that cause mist to disconnect, and release
			// nanobox
			if msg.Tags[0] == "err" {
				fmt.Printf(stylish.Warning("An unexpected error caused the sync stream to disconnect. Your sync should continue as normal, and you can see the log output from your dashboard."))
				break stream
			}

			//
			if err := json.Unmarshal([]byte(msg.Data), &log); err != nil {
				ui.LogFatal("[commands sync] json.Unmarshal() failed", err)
			}

			processLog(log)
		}

		// load historical logs
	} else {

		v := url.Values{}

		v.Add("level", fLevel)
		v.Add("reset", fmt.Sprintf("%v", fCount))

		logs := []string{}

		if err := api.DoRawRequest(nil, "GET", fmt.Sprintf("%v:6362/app?%v", config.Nanofile.IP, v.Encode()), &logs, nil); err != nil {
			ui.LogFatal("[commands sync] api.DoRawRequest() failed ", err)
		}

		ui.CPrint("[yellow]Showing last %v log entries...[reset]", strconv.Itoa(fCount))

		// display logs
		for _, log := range logs {
			fmt.Println("HERE??", log)
		}

	}
}

// processLog takes a Logvac or Stormpack log and breaks it apart into pieces that
// are then reconstructed in a 'digestible' way, colorized, and output to the
// terminal
func processLog(log Log) {

	time := time.Unix(0, int64(log.Time)*1000).Format(time.RFC822)

	//
	reFindLog := regexp.MustCompile(`^(\w+)\.(\S+)\s+(.*)$`)

	// example data stream parsed with regex
	// log.Log  = web1.apache[access] 69.92.84.90 - - [03/Dec/2013:19:59:57 +0000] \"GET / HTTP/1.1\" 200 183 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36\"\n"
	// subMatch[1] = web1
	// subMatch[2] = apache[access]
	// subMatch[3] = 69.92.84.90 - - [03/Dec/2013:19:59:57 +0000] \"GET / HTTP/1.1\" 200 183 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36\"\n"

	//
	config.Console.Debug("[commands.app_log.processLog] Raw log -> %#v", log)

	subMatch := reFindLog.FindStringSubmatch(log.Log)

	// ensure a subMatch and ensure subMatch has a length of 4, since thats how many
	// matches we're expecting
	if subMatch != nil && len(subMatch) >= 4 {

		service := subMatch[1]
		process := subMatch[2]
		entry := subMatch[3]

		//
		config.Console.Debug("[commands.app_log.processLog] Processed log -> service: %v, process: %v, entry: %v\n", service, process, entry)

		if _, ok := logProcesses[process]; !ok {
			logProcesses[process] = logColors[len(logProcesses)%len(logColors)]
		}

		ui.CPrint("[%v]%v - %v.%v :: %v[reset]", logProcesses[process], time, service, process, entry)

		// if we don't have a subMatch or its length is less than 4, just print w/e
		// is in the log
	} else {
		//
		config.Console.Debug("[commands.app_log.processLog] No submatches found -> %v - %v", time, log.Log)

		ui.CPrint("[light_red]%v - %v[reset]", time, log.Log)
	}

}
