package commands

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	pagodaAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

type (

	// AppLogCommand satisfies the Command interface for obtaining an app's historical
	// and streaming logs
	AppLogCommand struct{}

	// Logvac represents the system used for obtianing an app's historical logs
	Logvac struct {
		Logs  Logs
		Token string `json:"logvac_token"`
	}

	// Stormpack represents the system used for obtianing an app's streaming logs
	Stormpack struct {
		Command string   `json:"command"`
		Data    string   `json:"data"`
		Error   string   `json:"error"`
		Filters []string `json:"filters"`
		Keys    []string `json:"keys"`
		Success bool     `json:"success"`
		Token   string   `json:"stormpack_token"`
	}

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
	// "light_red",
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

// Help prints detailed help text for the app log command
func (c *AppLogCommand) Help() {
	ui.CPrintln(`
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
  `)
}

// Run attempts to display an app's logs. It takes count flag to designate how
// many logs to print, and a stream flag to indicate the live stream rather than
// historical. Logs are prased, colorized and printed to the terminal
func (c *AppLogCommand) Run(fApp string, opts []string, api *pagodaAPI.Client) {

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fCount int
	flags.IntVar(&fCount, "c", 100, "")
	flags.IntVar(&fCount, "count", 100, "")

	var fStream bool
	flags.BoolVar(&fStream, "l", false, "")
	flags.BoolVar(&fStream, "live", false, "")

	if err := flags.Parse(opts); err != nil {
		fmt.Println("Unable to parse flags. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda app:log", err)
	}

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	//
	app, err := api.GetApp(fApp)
	if err != nil {
		fmt.Printf("Oops! We could not find '%s'", fApp)
		os.Exit(1)
	}

	// if stream is true, we connect to the live logs
	if fStream {
		var stormpack Stormpack

		// authenticate stormpack
		fmt.Println("Authenticating... ")
		if err := api.DoRawRequest(&stormpack, "POST", "https://api.pagodabox.io/v1/auth/stormpack?auth_token="+api.AuthToken, nil, nil); err != nil {
			fmt.Println("Unable to complete API request. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda app:log", err)
		}

		// connect websocket
		fmt.Println("Connecting to live stream... ")
		url := "wss://smoke.pagodabox.io:443/subscribe/websocket?token=" + stormpack.Token
		ws, err := websocket.Dial(url, "", "http://localhost/")
		if err != nil {
			fmt.Println("Unable to establish websocket. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda app:log", err)
		}

		// subscribe to logs
		logs := Stormpack{Command: "subscribe", Filters: []string{app.ID, "log", "-member"}}
		if err := websocket.JSON.Send(ws, logs); err != nil {
			fmt.Println("Unable to send over websocket. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda app:log", err)
		}

		//
		var log Log

		// keep the socket open and listen for responses
		for {
			websocket.JSON.Receive(ws, &logs)

			if logs.Data != "" {
				if err := json.Unmarshal([]byte(logs.Data), &log); err != nil {
					fmt.Println("Unable to parse data. See ~/.pagodabox/log.txt for details")
					ui.Error("pagoda app:log", err)
				}

				processLog(log)
			}
		}

		// load historical logs
	} else {
		var logvac Logvac

		// authenticate logvac
		if err := api.DoRawRequest(&logvac, "POST", "https://api.pagodabox.io/v1/auth/logvac?auth_token="+api.AuthToken, nil, nil); err != nil {
			fmt.Println("Unable to complete API request. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda app:log", err)
		}

		// request historical logs
		if err := api.DoRawRequest(&logvac.Logs, "GET", "https://log.pagodabox.io/app/"+app.ID+"?limit="+strconv.Itoa(fCount), nil, map[string]string{"X-AUTH-TOKEN": logvac.Token}); err != nil {
			fmt.Println("Unable to complete API request. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda app:log", err)
		}

		// sort logs
		sort.Sort(logvac.Logs)

		ui.CPrintln("[yellow]Showing last " + strconv.Itoa(fCount) + " log entries...[reset]")

		// display logs
		for _, log := range logvac.Logs {
			processLog(log)
		}

	}
}

//
var reFindLog = regexp.MustCompile(`^(\w+)\.(\w+\[\w+\])\s+(.+)`)

// processLog takes a Logvac or Stormpack log and breaks it apart into pieces that
// are then reconstructed in a 'digestible' way, colorized, and output to the
// terminal
func processLog(log Log) {

	time := time.Unix(0, int64(log.Time)*1000).Format(time.RFC822)

	// example data stream parsed with regex
	// log.Log  = web1.apache[access] 69.92.84.90 - - [03/Dec/2013:19:59:57 +0000] \"GET / HTTP/1.1\" 200 183 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36\"\n"
	// subMatch[1] = web1
	// subMatch[2] = apache[access]
	// subMatch[3] = 69.92.84.90 - - [03/Dec/2013:19:59:57 +0000] \"GET / HTTP/1.1\" 200 183 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36\"\n"

	subMatch := reFindLog.FindStringSubmatch(log.Log)
	// if subMatch == nil { }

	// check for a length of 4 because thats how many matches we're expecting
	if len(subMatch) >= 4 {

		service := subMatch[1]
		process := subMatch[2]
		entry := subMatch[3]

		if _, ok := logProcesses[process]; !ok {
			logProcesses[process] = logColors[len(logProcesses)%len(logColors)]
		}

		ui.CPrintln("[" + logProcesses[process] + "]" + time + " - " + service + "." + process + " ::" + entry + "[reset]")

		// if we don't have a len of 4 we probably have 0, but just incase we'll print
		// whatever we get.
	} else {
		ui.CPrintln("[light_red]" + time + log.Log + "[reset]")
	}

}
