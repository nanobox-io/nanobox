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
	"regexp"
	// "sort"
	// "strconv"
	"time"

	mist "github.com/pagodabox/golang-mist"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

type (

	// LogCommand satisfies the Command interface for obtaining an app's historical
	// and streaming logs
	LogCommand struct{}

	// Logs represents a slice of Log's
	Logs []Log

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

// functions for sorting logs by timestamp
func (l Logs) Len() int {
	return len(l)
}

func (l Logs) Less(i, j int) bool {
	return l[i].Time < l[j].Time
}

func (l Logs) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Help
func (c *LogCommand) Help() {
	ui.CPrint(`
Description:
  Provides log output for an application.

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

  If [count] is not provided, will show the last 100 lines of the log.

  If [live] is not provided, will default to showing the last 100 lines.

Usage:
  pagoda log [-a app-name] [-c count] [-l]
  pagoda app:log [-a app-name] [-c count] [-l]

  ex. pagoda log -a app-name -c 100 -l

Options:
  -a, --app [app-name]
    The name of the app you want to view logs for.

  -c, --count [count]
    The number of lines of the log you wish to view.

  -l, --live
    Enable live stream
    // emergency, alert, critical, error, warning, notice, informational, debug, log

  --level
    Log level
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
	flags.BoolVar(&fStream, "l", false, "")
	flags.BoolVar(&fStream, "live", false, "")

	var fLevel string
	flags.StringVar(&fLevel, "level", "info", "")

	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.app_log] flags.Parse()", err)
	}

	// if stream is true, we connect to the live logs
	if fStream {

		// connect websocket
		fmt.Println("Connecting to live stream...")

		// subscribe to mist
		client := mist.Client{Host: config.Nanofile.IP, Port: "1445"}
		if err := client.Connect(); err != nil {
			ui.LogFatal("[commands deploy] client.Connect() failed ", err)
		}

		defer client.Close()

		//
		sub, err := client.Subscribe([]string{"app", fLevel})
		if err != nil {
			config.Console.Warn("Failed to subscribe to 'mist' updates... %v", err)
		}

		//
		config.Console.Debug("[commands.app_log.Run] Subscribing to logs: %#v", sub)

		//
		var log Log

		//
		fmt.Println("Waiting for data...")

		//
		for msg := range client.Data {

			data := make(map[string]string)

			if data["log"] != "" {

				//
				if err := json.Unmarshal([]byte(msg.Data), &log); err != nil {
					ui.LogFatal("[commands.app_log] json.Unmarshal() failed", err)
				}

				processLog(log)
			}
		}

		// load historical logs
	} else {

		// logsURL := fmt.Sprintf("https://log.pagodabox.io/app/%v?limit=%v", safeID, strconv.Itoa(fCount))

		// //
		// config.Console.Debug("[commands.app_log.Run] Requesting historical logs from: %v", logsURL)

		// // request historical logs
		// if err := api.DoRawRequest(&logvac.Logs, "GET", logsURL, nil, map[string]string{"X-AUTH-TOKEN": logvac.Token}); err != nil {
		//   ui.LogFatal("[commands.app_log] api.DoRawRequest() failed", err)
		// }

		// // sort logs
		// sort.Sort(logvac.Logs)

		// ui.CPrint("[yellow]Showing last %v log entries...[reset]", strconv.Itoa(fCount))

		// // display logs
		// for _, log := range logvac.Logs {
		//   processLog(log)
		// }

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
