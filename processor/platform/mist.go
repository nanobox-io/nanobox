package platform

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanopack/mist/clients"
)

// processMistListen ...
type processMistListen struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("mist_log", mistListenFn)
}

//
func mistListenFn(control processor.ProcessControl) (processor.Processor, error) {
	return processMistListen{control}, nil
}

//
func (mistListen processMistListen) Results() processor.ProcessControl {
	return mistListen.control
}

//
func (mistListen processMistListen) Process() error {
	mist := models.Service{}
	bucket := fmt.Sprintf("%s_%s", config.AppID(), mistListen.control.Env)
	data.Get(bucket, "mist", &mist)

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
