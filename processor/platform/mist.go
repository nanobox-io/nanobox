package platform

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/nanopack/mist/clients"

	"github.com/nanobox-io/nanobox/models"
)

// MistListen ...
type MistListen struct {
	App models.App
}

//
func (mistListen MistListen) Run() error {
	mist, err := models.FindComponentBySlug(mistListen.App.ID, "mist")

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

	// loop waiting for messages or signals
	// if we recieve a kill signal quit
	// messages will be displayed
	for {
		select {
		case msg := <-client.Messages():
			fmt.Printf("message: %+v\n", msg)
		case <-sigChan:
			return nil
		}
	}
}
