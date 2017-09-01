// Package log is a processor for live streaming and pulling historic production logs.
package log

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanopack/logvac/core"
	"github.com/nanopack/mist/core"
	"golang.org/x/net/websocket"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/odin"
)

// Tail tails production logs for an app.
func Tail(envModel *models.Env, app string, logOpts models.LogOpts) error {
	appID := app

	// fetch the remote
	remote, ok := envModel.Remotes[appID]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(remote.Endpoint)
		// set the app id
		appID = remote.ID
	}

	// set odins endpoint if the arguement is passed
	if endpoint := registry.GetString("endpoint"); endpoint != "" {
		odin.SetEndpoint(endpoint)
	}

	// todo: don't assume app name, just message and die
	// set the app id to the directory name if it's default
	if appID == "default" {
		appID = config.AppName()
	}

	// validate access to the app
	if err := helpers.ValidateOdinApp(appID); err != nil {
		return util.ErrorAppend(err, "unable to validate app")
	}

	mistConfig, err := getMistConfig(envModel, appID)
	if err != nil {
		return util.ErrorAppend(err, "unable to generate mist config")
	}

	// fmt.Println("mistConfig", mistConfig)
	err = mistListen(mistConfig.Token, mistConfig.URL, logOpts)
	if err != nil {
		return util.ErrorAppend(err, "failed to subscribe to logs")
	}

	return nil
}

// Print prints historic production logs for an app.
func Print(envModel *models.Env, app string, logOpts models.LogOpts) error {
	appID := app

	// fetch the remote
	remote, ok := envModel.Remotes[appID]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(remote.Endpoint)
		// set the app id
		appID = remote.ID
	}

	// set odins endpoint if the arguement is passed
	if endpoint := registry.GetString("endpoint"); endpoint != "" {
		odin.SetEndpoint(endpoint)
	}

	// todo: don't assume app name, just message and die
	// set the app id to the directory name if it's default
	if appID == "default" {
		appID = config.AppName()
	}

	// validate access to the app
	if err := helpers.ValidateOdinApp(appID); err != nil {
		return util.ErrorAppend(err, "unable to validate app")
	}

	token, url, err := odin.GetComponent(appID, "logger")
	if err != nil {
		lumber.Error("deploy:setMistToken:GetMist(%s): %s", appID, err.Error())
		err = util.ErrorAppend(err, "failed to fetch logvac information from nanobox")
		return err
	}

	err = fetchLogs(token, url, logOpts)
	if err != nil {
		return util.ErrorAppend(err, "failed to fetch logs")
	}

	return nil
}

type MistConfig struct {
	URL   string
	Token string
}

func getMistConfig(envModel *models.Env, appID string) (*MistConfig, error) {

	token, url, err := odin.GetComponent(appID, "pusher")
	if err != nil {
		lumber.Error("deploy:setMistToken:GetMist(%s): %s", appID, err.Error())
		err = util.ErrorAppend(err, "failed to fetch mist information from nanobox")
		return nil, err
	}

	return &MistConfig{url, token}, nil
}

// mistListen will subscribe to mist and print incoming logs.
func mistListen(token, url string, logOpts models.LogOpts) error {
	logFollow := logOpts.Follow
	// connect to the mist server
	var wsConn *websocket.Conn
	clientConnect := func() (err error) {
		wsConn, err = newMistClient(token, url)
		return err
	}
	if err := util.Retry(clientConnect, 3, time.Second); err != nil {
		return err
	}

	// subscribe to all logs
	if err := subscribe(wsConn); err != nil {
		return err
	}

	// catch kill signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// if `-f` wasn't explicity called, print what we are doing
	if !logFollow {
		fmt.Printf(`
Connected to streaming logs:
ctrl + c to quit
------------------------------------------------
waiting for output...

`)
	}

	// loop waiting for messages or signals if we recieve a kill signal quit
	// messages will be displayed
	// msgChan := client.Messages()
	for {
		select {
		case msg := <-messageChan:
			display.FormatLogMessage(msg, logOpts.Raw)
		case <-sigChan:
			return nil
		}
	}
}

// todo: make part of mist wsclient

var messageChan chan mist.Message

func newMistClient(token, address string) (*websocket.Conn, error) {
	origin := "https://nanoapp.localhost"
	url := "wss://" + address + ":1446/subscribe/websocket?X-AUTH-TOKEN=" + token

	config, err := websocket.NewConfig(url, origin)
	if err != nil {
		return nil, fmt.Errorf("failed to create config - %s", err.Error())
	}

	config.TlsConfig = &tls.Config{InsecureSkipVerify: true}

	ws, err := websocket.DialConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial mist - %s", err.Error())
	}

	messageChan = make(chan mist.Message, 1)
	// connection loop (blocking); continually read off the connection. Once something
	// is read, check to see if it's a message the client understands to be one of
	// its commands. If so attempt to execute the command.
	go func() {
		decoder := json.NewDecoder(ws)

		for decoder.More() {
			msg := mist.Message{}

			// decode an array value (Message)
			if err := decoder.Decode(&msg); err != nil {
				// invalid character '\x15' looking for beginning of value
				if strings.Contains(err.Error(), "invalid character '\\x15'") {
					fmt.Printf("Must dial TLS - %s\n", err.Error())
					return
				}

				// an error decoding should be sent to the user
				reader := decoder.Buffered()
				bytes, _ := ioutil.ReadAll(reader)
				msg.Error = string(bytes)
			}

			messageChan <- msg
		}
	}()

	return ws, nil
}

type mistCommand struct {
	Command string   `json:"command"`
	Tags    []string `json:"tags"`
}

func subscribe(ws *websocket.Conn) error {
	b, err := json.Marshal(mistCommand{"subscribe", []string{"log"}})
	if err != nil {
		return err
	}

	_, err = ws.Write(b)
	return err
}

// fetchLogs fetches and prints historic logs
func fetchLogs(token, url string, logOpts models.LogOpts) error {
	numLogs := logOpts.Number
	body, err := rest(url, "GET", fmt.Sprintf("/logs?type=app&id=&start=0&limit=%d", numLogs), token, token)
	if err != nil {
		return fmt.Errorf("failed to get logs - %s", err.Error())
	}

	msgs := []logvac.Message{}
	err = json.Unmarshal(body, &msgs)
	if err != nil {
		return util.Err{
			Message: fmt.Sprintf("Failed to process historic logs - %s", err.Error()),
			Suggest: "Please upgrade your logging component and try again.",
		}
	}

	for i := range msgs {
		display.FormatLogvacMessage(msgs[i], logOpts.Raw)
	}

	return nil
}

func rest(ip, method, route, auth, user string) ([]byte, error) {
	req, _ := http.NewRequest(method, fmt.Sprintf("https://%s:6361%s", ip, route), nil)
	req.Header.Add("X-AUTH-TOKEN", auth)
	req.Header.Add("X-USER-TOKEN", user)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to %s %s - %s", method, route, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Status '200' expected, got '%d'", res.StatusCode)
	}

	b, _ := ioutil.ReadAll(res.Body)

	return b, nil
}
