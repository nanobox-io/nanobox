// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	// "time"

	"github.com/spf13/cobra"

	"github.com/pagodabox/golang-mist"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	// "github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Provides the last 100 lines of historical log output by default.",
	Long: `
Description:
  Provides the last 100 lines of historical log output by default.`,

	Run: nanoLog,
}

// Log represents the structure of a log returned from Logvac or Stormpack
type Log struct {
	Level string `json:"level"`
	Log   string `json:"log"`
	Time  string `json:"time"`
}

var (
	// a map of each type of 'process' that we encounter to then be used when assigning
	// a unique color to that 'process'
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

//
func init() {
	logCmd.Flags().IntVarP(&fCount, "count", "c", 100, "Specifies the number of lines to output from the historical log.")
	logCmd.Flags().BoolVarP(&fStream, "stream", "s", false, "Streams logs live")
	logCmd.Flags().StringVarP(&fLevel, "level", "l", "info", "Filters logs by one of the following levels: debug > info > warn > error > fatal")
}

// nanoLog
func nanoLog(ccmd *cobra.Command, args []string) {

	// if stream is true, we connect to the live logs
	if fStream {
		fmt.Printf(stylish.Bullet("Connecting to live stream..."))

		// connect 'mist' to the server running on the guest machine
		client, err := mist.NewRemoteClient(config.Nanofile.IP + ":1445")
		if err != nil {
			ui.LogFatal("[commands/log] client.Connect() failed ", err)
		}
		defer client.Close()

		// subscribe to 'log' updates
		logTags := []string{"log", fLevel}
		if err := client.Subscribe([]string{"log", fLevel}); err != nil {
			fmt.Printf(stylish.Warning("Nanobox failed to subscribe to app logs."))
		}
		defer client.Unsubscribe(logTags)

		//
		fmt.Printf(stylish.Bullet("Connecting to live stream..."))

	stream:
		for msg := range client.Messages() {

			//
			log := Log{}

			// check for any error message that cause mist to disconnect, and release
			// nanobox
			if msg.Tags[0] == "err" {
				fmt.Printf(stylish.Warning("An unexpected error caused the sync stream to disconnect."))
				break stream
			}

			//
			if err := json.Unmarshal([]byte(msg.Data), &log); err != nil {
				ui.LogFatal("[commands/log] json.Unmarshal() failed", err)
			}

			processLog(log)
		}

		// load historical logs
	} else {
		fmt.Printf(stylish.Bullet(fmt.Sprintf("Showing last %v entries:", fCount)))

		//
		v := url.Values{}
		v.Add("level", fLevel)
		v.Add("limit", fmt.Sprintf("%v", fCount))

		res, err := http.Get(fmt.Sprintf("http://%v:6362/app?%v", config.Nanofile.IP, v.Encode()))
		if err != nil {
			ui.LogFatal("[commands/log] http.Get() failed", err)
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
					ui.LogFatal("[commands/log] bufio.ReadBytes() failed", err)
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
	config.Console.Debug("[commands/log] Raw log -> %#v", log)

	subMatch := reFindLog.FindStringSubmatch(log.Log)

	// ensure a subMatch and ensure subMatch has a length of 4, since thats how many
	// matches we're expecting
	if subMatch != nil && len(subMatch) >= 4 {

		service := subMatch[1]
		process := subMatch[2]
		entry := subMatch[3]

		//
		config.Console.Debug("[commands/log] Processed log -> service: %v, process: %v, entry: %v\n", service, process, entry)

		if _, ok := logProcesses[process]; !ok {
			logProcesses[process] = logColors[len(logProcesses)%len(logColors)]
		}

		ui.CPrint("[%v]%v - %v.%v :: %v[reset]", logProcesses[process], log.Time, service, process, entry)

		// if we don't have a subMatch or its length is less than 4, just print w/e
		// is in the log
	} else {
		//
		config.Console.Debug("[commands/log] No submatches found -> %v - %v", log.Time, log.Log)

		ui.CPrint("[light_red]%v - %v[reset]", log.Time, log.Log)
	}

}
