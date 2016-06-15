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
	// confirm the provider is an accessable one that we support.
	return processMistListen{control}, nil
}

//
func (mistListen processMistListen) Results() processor.ProcessControl {
	return mistListen.control
}

//
func (mistListen processMistListen) Process() error {
	mist := models.Service{}
	data.Get(config.AppName(), "mist", &mist)

	//
	client, err := clients.New(mist.ExternalIP+":1445", "123")
	if err != nil {
		return err
	}

	//
	if err := client.Subscribe([]string{"log"}); err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	//
	for {
		select {
		case msg := <-client.Messages():
			fmt.Printf("message: %+v\n", msg)
		case <-sigChan:
			return nil
		}
	}
}
