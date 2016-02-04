//
package dev

import (
	"fmt"
	"net/url"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/server"
	mistutil "github.com/nanobox-io/nanobox/util/server/mist"
)

var (

	//
	logCmd = &cobra.Command{
		Hidden: true,

		Use:   "log",
		Short: "Provides the last 100 lines of historical log output (default).",
		Long:  ``,

		Run: log,
	}

	//
	count  int    // number of logs to show
	level  string // log level of logs to show
	offset int    // log entry to begin showing logs from
	stream bool   // whether to stream the logs or not
)

//
func init() {
	logCmd.Flags().BoolVarP(&stream, "stream", "s", false, "Streams logs live")
	logCmd.Flags().IntVarP(&count, "count", "c", 100, "Specifies the number of lines to output from the historical log.")
	logCmd.Flags().StringVarP(&level, "level", "l", "info", "Filters logs by one of the following levels: debug > info > warn > error > fatal")
	logCmd.Flags().IntVarP(&offset, "offset", "o", 0, "Specifies the entry at which to start pulling <count> from")
}

// log
func log(ccmd *cobra.Command, args []string) {

	// if stream is true, we connect to the live logs...
	if stream {

		fmt.Printf(stylish.Bullet("Opening log stream"))

		// stream logs (blocking)
		mistutil.Stream([]string{"log", "app"}, mistutil.ProcessLogStream)

		// ...otherwise load historical logs
	} else {

		//
		v := url.Values{}
		v.Add("level", level)
		v.Add("limit", fmt.Sprintf("%v", count))
		v.Add("offset", fmt.Sprintf("%v", offset))

		// show Mist history
		if err := server.Logs(v.Encode()); err != nil {
			server.Fatal("[commands/log] server.Logs() failed", err.Error())
		}
	}
}
