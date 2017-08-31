package display

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/mitchellh/colorstring"
	"github.com/nanopack/logvac/core"
	"github.com/nanopack/mist/core"
)

// Entry represents the data comming back from a mist message (mist.Message.Data)
type Entry struct {
	Time     time.Time `json:"time"`     // "2016-09-07T20:33:34.446275741Z"
	UTime    int       `json:"utime"`    // 1473280414446275741
	ID       string    `json:"id"`       // "mist"
	Tag      []string  `json:"tag"`      // ["mist[daemon]", "mist"]
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
func FormatLogMessage(msg mist.Message) {

	// set the time output format
	layout := "Mon Jan 02 15:04:05 2006" // time.RFC822

	// unmarshal the message data as an Entry
	entry := Entry{}
	if err := json.Unmarshal([]byte(msg.Data), &entry); err != nil {
		message := fmt.Sprintf("[light_red]%s :: %s\n[reset]%s", time.Now().Format(layout), msg.Data, fmt.Sprintf("Failed to process entry - '%s'. Please upgrade your logging component and try again.", err.Error()))
		fmt.Println(colorstring.Color(message))
		return
	}

	//
	fmtMsg := regexp.MustCompile(`\s?\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d+Z|\s?\d{4}-\d{2}-\d{2}[_T]\d{2}:\d{2}:\d{2}.\d{5}|\s?\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}|\s?\[\d{2}\/\w{3}\/\d{4}\s\d{2}:\d{2}:\d{2}\]?`).ReplaceAllString(entry.Message, "")

	// set default (shouldn't ever be needed)
	entryTag := "localApp"
	if len(entry.Tag) != 0 {
		entryTag = entry.Tag[0]
	}

	// for each new entryTag assign it a color to be used when output
	if _, ok := logProcesses[entryTag]; !ok {
		logProcesses[entryTag] = logColors[len(logProcesses)%len(logColors)]
	}

	// return our pretty entry
	message := fmt.Sprintf("[%s]%s %s (%s) :: %s[reset]", logProcesses[entryTag], fmt.Sprintf(entry.Time.Format(layout)), entry.ID, entryTag, fmtMsg)
	fmt.Println(colorstring.Color(message))
	return
}

// FormatLogvacMessage takes a Logvac/Mist and formats it into a pretty message to be
// output to the terminal
func FormatLogvacMessage(msg logvac.Message) {

	// set the time output format
	layout := "Mon Jan 02 15:04:05 2006" // time.RFC822

	// unmarshal the message data as an Entry
	entry := Entry{
		Time:     msg.Time,
		UTime:    int(msg.UTime),
		ID:       msg.Id,
		Tag:      msg.Tag,
		Type:     msg.Type,
		Priority: msg.Priority,
		Message:  msg.Content,
	}

	fmtMsg := regexp.MustCompile(`\s?\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d+Z|\s?\d{4}-\d{2}-\d{2}[_T]\d{2}:\d{2}:\d{2}.\d{5}|\s?\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}|\s?\[\d{2}\/\w{3}\/\d{4}\s\d{2}:\d{2}:\d{2}\]?`).ReplaceAllString(entry.Message, "")

	// set default (shouldn't ever be needed)
	entryTag := "localApp"
	if len(entry.Tag) != 0 {
		entryTag = entry.Tag[0]
	}

	// for each new entryTag assign it a color to be used when output
	if _, ok := logProcesses[entryTag]; !ok {
		logProcesses[entryTag] = logColors[len(logProcesses)%len(logColors)]
	}

	// return our pretty entry
	message := fmt.Sprintf("[%s]%s %s (%s) :: %s[reset]", logProcesses[entryTag], fmt.Sprintf(entry.Time.Format(layout)), entry.ID, entryTag, fmtMsg)
	fmt.Println(colorstring.Color(message))
	return
}
