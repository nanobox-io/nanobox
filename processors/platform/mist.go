package platform

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nanopack/mist/clients"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
)

// MistListen ...
func MistListen(appModel *models.App) error {
	mist, err := models.FindComponentBySlug(appModel.ID, "mist")
	if err != nil {
		return err
	}

	// connect to the mist server
	var client *clients.TCP
	clientConnect := func() (err error) {
		client, err = clients.New(mist.IPAddr()+":1445", "123")
		return err
	}
	if err := util.Retry(clientConnect, 3, time.Second); err != nil {
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
ctrl + c to quit
------------------------------------------------
waiting for output...

`)

	// loop waiting for messages or signals if we recieve a kill signal quit
	// messages will be displayed
	for {
		select {
		case msg := <-client.Messages():
			display.FormatLogMessage(msg, false)
		case <-sigChan:
			return nil
		}
	}
}
