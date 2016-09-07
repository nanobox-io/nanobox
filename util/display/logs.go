package display

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/nanopack/mist/core"
)

// Entry represents the data comming back from a mist message (mist.Message.Data)
type Entry struct {
	Time     time.Time `json:"time"`     // "2016-09-07T20:33:34.446275741Z"
	UTime    int       `json:"utime"`    // 1473280414446275741
	ID       string    `json:"id"`       // "mist"
	Tag      string    `json:"tag"`      // "mist[daemon]"
	Type     string    `json:"type"`     // "app"
	Priority int       `json:"priority"` // 4
	Message  string    `json:"message"`  // "2016-09-07T20:33:34.44586 2016-09-07 20:33:34 INFO  Api Listening on https://0.0.0.0:6361..."
}

var (
	// a map of each type of 'process' that we encounter to then be used when
	// assigning a unique color to that 'process'
	logProcesses = make(map[string]string)

	// an array of the colors used to colorize the logs
	logColors = [10]string{"green", "yellow", "blue", "magenta", "cyan", "light_green", "light_yellow", "light_blue", "light_magenta", "light_cyan"}
)

// FormatLogMessage takes a Logvac/Mist and formats it into a pretty message to be
// output to the terminal
func FormatLogMessage(msg mist.Message) string {

	// set the time output format
	layout := "Mon Jan 02 15:04:05 2006" // time.RFC822

	// unmarshal the message data as an Entry
	entry := Entry{}
	if err := json.Unmarshal([]byte(msg.Data), &entry); err != nil {
		return fmt.Sprintf("[light_red]%s :: %s[reset]", time.Now().Format(layout), "Failed to process entry...")
	}

	//
	fmtMsg := regexp.MustCompile(`\s?\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d+Z|\s?\d{4}-\d{2}-\d{2}[_T]\d{2}:\d{2}:\d{2}.\d{5}|\s?\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}|\s?\[\d{2}\/\w{3}\/\d{4}\s\d{2}:\d{2}:\d{2}\]?`).ReplaceAllString(entry.Message, "")

	// for each new entry.Tag assign it a color to be used when output
	if _, ok := logProcesses[entry.Tag]; !ok {
		logProcesses[entry.Tag] = logColors[len(logProcesses)%len(logColors)]
	}

	// return our pretty entry
	return fmt.Sprintf("[%s]%s %s (%s) :: %s[reset]", logProcesses[entry.Tag], fmt.Sprintf(entry.Time.Format(layout)), entry.ID, entry.Tag, fmtMsg)
}
