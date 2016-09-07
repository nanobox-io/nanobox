package platform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/nanopack/mist/clients"
	"github.com/nanopack/mist/core"
	printutil "github.com/sdomino/go-util/print"

	"github.com/nanobox-io/nanobox/models"
)

// Log represents the data comming back from a mist message (mist.Message.Data)
type Log struct {
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

// MistListen ...
func MistListen(appModel *models.App) error {
	mist, err := models.FindComponentBySlug(appModel.ID, "mist")

	// connect to the mist server
	client, err := clients.New(mist.ExternalIP+":1445", "123")
	if err != nil {
		return err
	}

	// subscribe to all logs
	if err := client.Subscribe([]string{"log"}); err != nil {
		return err
	}

	// catch kill signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	fmt.Printf(`
Connected to streaming logs:
------------------------------------------------
waiting for output...

`)

	// loop waiting for messages or signals if we recieve a kill signal quit
	// messages will be displayed
	for {
		select {
		case msg := <-client.Messages():
			printutil.Color(formatMessage(msg))
		case <-sigChan:
			return nil
		}
	}
}

// formatMessage takes a Logvac/Mist and formats it into a pretty message to be
// output to the terminal
func formatMessage(msg mist.Message) string {

	// set the time output format
	layout := "Mon Jan 02 15:04:05 2006" // time.RFC822

	// unmarshal the message data as a Log
	log := Log{}
	if err := json.Unmarshal([]byte(msg.Data), &log); err != nil {
		return fmt.Sprintf("[light_red]%s :: %s[reset]", time.Now().Format(layout), "Failed to process log...")
	}

	//
	shortDateTime := fmt.Sprintf(log.Time.Format(layout))
	entry := regexp.MustCompile(`\s?\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d+Z|\s?\d{4}-\d{2}-\d{2}[_T]\d{2}:\d{2}:\d{2}.\d{5}|\s?\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}|\s?\[\d{2}\/\w{3}\/\d{4}\s\d{2}:\d{2}:\d{2}\]?`).ReplaceAllString(log.Message, "")

	// for each new log.Tag assign it a color to be used when output
	if _, ok := logProcesses[log.Tag]; !ok {
		logProcesses[log.Tag] = logColors[len(logProcesses)%len(logColors)]
	}

	// return our pretty log
	return fmt.Sprintf("[%s]%s %s (%s) :: %s[reset]", logProcesses[log.Tag], shortDateTime, log.ID, log.Tag, entry)
}
