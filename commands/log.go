// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	// "time"

	"github.com/spf13/cobra"

	"github.com/pagodabox/golang-mist"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var logCmd = &cobra.Command{
	Hidden: true,

	Use:   "log",
	Short: "Provides the last 100 lines of historical log output by default.",
	Long: `
Description:
  Provides the last 100 lines of historical log output by default.`,

	PreRun: vmIsRunning,
	Run:    nanoLog,
}

// Log represents the structure of a log returned from Logtap or Stormpack
type Log struct {
	Content  string `json:"content"`
	Priority int    `json:"priority"`
	Time     string `json:"time"`
	Type     string `json:"Type"`
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
	logCmd.Flags().BoolVarP(&fStream, "stream", "s", false, "Streams logs live")

	logCmd.Flags().IntVarP(&fCount, "count", "c", 100, "Specifies the number of lines to output from the historical log.")
	logCmd.Flags().StringVarP(&fLevel, "level", "l", "info", "Filters logs by one of the following levels: debug > info > warn > error > fatal")
	logCmd.Flags().IntVarP(&fOffset, "offset", "o", 0, "Specifies the entry at which to start pulling <count> from")
}

// nanoLog
func nanoLog(ccmd *cobra.Command, args []string) {

	// if stream is true, we connect to the live logs...
	if fStream {

		// connect 'mist' to the server running on the guest machine
		client, err := mist.NewRemoteClient(config.MistURI)
		if err != nil {
			config.Fatal("[commands/log] client.Connect() failed ", err)
		}
		defer client.Close()

		// subscribe to 'log' updates
		logTags := []string{"log", fLevel}
		if err := client.Subscribe([]string{"log", fLevel}); err != nil {
			fmt.Printf(stylish.Warning("Nanobox failed to subscribe to app logs."))
		}
		defer client.Unsubscribe(logTags)

		//
		fmt.Printf("Streaming App Logs:\n")
		for msg := range client.Messages() {

			//
			log := Log{}
			if err := json.Unmarshal([]byte(msg.Data), &log); err != nil {
				config.Fatal("[commands/log] json.Unmarshal() failed", err)
			}

			//
			processLog(log)
		}

		// ...otherwise load historical logs
	} else {

		//
		v := url.Values{}
		v.Add("level", fLevel)
		v.Add("limit", fmt.Sprintf("%v", fCount))
		v.Add("offset", fmt.Sprintf("%v", fOffset))

		//
		res, err := http.Get(fmt.Sprintf("http://%v/logs?%v", config.ServerURI, v.Encode()))
		if err != nil {
			config.Fatal("[commands/log] http.Get() failed", err)
		}
		defer res.Body.Close()

		//
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			config.Fatal("[commands/log] ioutil.ReadAll() failed", err)
		}

		//
		logs := []Log{}
		if err := json.Unmarshal(b, &logs); err != nil {
			config.Fatal("[commands/log] json.Unmarshal() failed", err)
		}

		//
		fmt.Printf(stylish.Bullet("Showing last %v entries:", len(logs)))
		for _, log := range logs {
			processLog(log)
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
