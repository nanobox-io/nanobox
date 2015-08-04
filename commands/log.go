// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	// "time"

	"github.com/pagodabox/golang-mist"
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
		Level string `json:"level"`
		Log   string `json:"log"`
		Time  string `json:"time"`
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
		fmt.Printf(stylish.Bullet(fmt.Sprintf("Showing last %v entries:", fCount)))

		//
		v := url.Values{}
		v.Add("level", fLevel)
		v.Add("reset", fmt.Sprintf("%v", fCount))

		res, err := http.Get(fmt.Sprintf("http://%v:6362/app?%v", config.Nanofile.IP, v.Encode()))
		if err != nil {
			ui.LogFatal("[commands.log] http.Get() failed", err)
		}
		defer res.Body.Close()

		//
		reParseLog := regexp.MustCompile(`\[(.*)\] \[(.*)\] (.*)`)

		// read response body, which should be our version string
		r := bufio.NewReader(res.Body)
		for {
			b, err := r.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					ui.LogFatal("[commands.log] bufio.ReadBytes() failed", err)
				}
			}

			//
			subMatch := reParseLog.FindStringSubmatch(string(b))

			// ensure a subMatch and ensure subMatch has a length of 4, since thats how many
			// matches we're expecting
			if subMatch != nil && len(subMatch) >= 4 {
				processLog(Log{Level: subMatch[2], Log: subMatch[3], Time: subMatch[1]})
			}
		}
	}
}

// processLog takes a Logvac or Stormpack log and breaks it apart into pieces that
// are then reconstructed in a 'digestible' way, colorized, and output to the
// terminal
func processLog(log Log) {

	// t := time.Now(log.Time).Format(time.RFC822)
	// t, err := time.Parse("01/02 03:04:05PM '06 -0700", log.Time)
	// if err != nil {
	// 	fmt.Println("TIME BONK!", err)
	// }

	//
	reFindLog := regexp.MustCompile(`^(\w+)\.(\S+)\s+(.*)$`)

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

		ui.CPrint("[%v]%v - %v.%v :: %v[reset]", logProcesses[process], log.Time, service, process, entry)

		// if we don't have a subMatch or its length is less than 4, just print w/e
		// is in the log
	} else {
		//
		config.Console.Debug("[commands.app_log.processLog] No submatches found -> %v - %v", log.Time, log.Log)

		ui.CPrint("[light_red]%v - %v[reset]", log.Time, log.Log)
	}

}
