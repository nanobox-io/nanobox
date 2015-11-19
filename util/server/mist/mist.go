//
package mist

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	mistClient "github.com/nanopack/mist/core"

	"github.com/nanobox-io/nanobox/config"
	printutil "github.com/nanobox-io/nanobox/util/print"
)

//
type (

	// Model
	Model struct {
		Action   string `json:"action"`
		Document struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"document"`
		Name string `json:"model"`
	}

	// Log
	Log struct {
		Content  string `json:"content"`
		Priority int    `json:"priority"`
		Time     string `json:"time"`
		Type     string `json:"type"`
	}
)

//
var (

	// subscriptions is a list of tags that have been used to subscribe with either
	// Listen or Stream; when creating a new Listner/Streamer if the tags have
	// already been used, it stops double subscription
	subscriptions = make(map[string]struct{})

	// a map of each type of 'process' that we encounter to then be used when
	// assigning a unique color to that 'process'
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

// Listen connects a to mist, subscribes tags, and listens for 'model' updates
func Listen(tags []string, handle func(string) error) error {

	// only subscribe if a subscription doesn't already exist
	if _, ok := subscriptions[strings.Join(tags, "")]; ok {
		return nil
	}

	// connect to mist
	client, err := mistClient.NewRemoteClient(config.MistURI)
	if err != nil {
		config.Fatal("[util/server/mist/mist] mist.NewRemoteClient() failed", err.Error())
	}
	defer client.Close()

	// this is a bandaid to fix a race condition in mist when immediatly subscribing
	// after connecting a client; once this is fixed in mist this can be removed
	<-time.After(time.Second * 1)

	// subscribe
	if err := client.Subscribe(tags); err != nil {
		config.Fatal("[util/server/mist/mist] client.Subscribe() failed", err.Error())
	}
	defer delete(subscriptions, strings.Join(tags, ""))

	// add tags to list of subscriptions
	subscriptions[strings.Join(tags, "")] = struct{}{}

	//
	model := Model{}
	for msg := range client.Messages() {

		// unmarshal the incoming Message
		if err := json.Unmarshal([]byte(msg.Data), &model); err != nil {
			config.Fatal("[util/server/mist/mist] json.Unmarshal() failed", err.Error())
		}

		// handle the status; when the handler returns false, it's time to break the
		// stream
		return handle(model.Document.Status)
	}

	return nil
}

// Stream connects to mist, subscribes tags, and logs Messages
func Stream(tags []string, handle func(Log)) {

	// add log level to tags
	tags = append(tags, config.LogLevel)

	// if this subscription already exists, exit; this prevents double subscriptions
	if _, ok := subscriptions[strings.Join(tags, "")]; ok {
		return
	}

	// connect to mist
	client, err := mistClient.NewRemoteClient(config.MistURI)
	if err != nil {
		config.Fatal("[util/server/mist/mist] mist.NewRemoteClient() failed", err.Error())
	}
	defer client.Close()

	// this is a bandaid to fix a race condition in mist when immediatly subscribing
	// after connecting a client; once this is fixed in mist this can be removed
	<-time.After(time.Second * 1)

	// subscribe
	if err := client.Subscribe(tags); err != nil {
		config.Fatal("[util/server/mist/mist] client.Subscribe() failed", err.Error())
	}
	defer delete(subscriptions, strings.Join(tags, ""))

	// add tags to list of subscriptions
	subscriptions[strings.Join(tags, "")] = struct{}{}

	//
	for msg := range client.Messages() {

		//
		log := Log{}

		// unmarshal the incoming Message
		if err := json.Unmarshal([]byte(msg.Data), &log); err != nil {
			config.Fatal("[util/server/mist/mist] json.Unmarshal() failed", err.Error())
		}

		//
		handle(log)
	}
}

// ProcessLog takes a Logvac or Stormpack log and breaks it apart into pieces that
// are then reconstructed in a 'digestible' way, colorized, and output to the
// terminal
func ProcessLog(log Log) {

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

		// print.Color("[%v]%v - %v.%v :: %v[reset]", logProcesses[process], log.Time, service, process, content)
		printutil.Color("[%v]%v (%v) :: %v[reset]", logProcesses[process], service, process, content)

		// if we don't have a subMatch or its length is less than 4, just print w/e
		// is in the log
	} else {
		printutil.Color("[light_red]%v - %v[reset]", log.Time, log.Content)
	}
}
