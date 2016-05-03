package nanopack

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanopack/mist/clients"
)

type mistListen struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("mist_log", mistListenFunc)
}

func mistListenFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return mistListen{config}, nil
}

func (self mistListen) Results() processor.ProcessConfig {
	return self.config
}

func (self mistListen) Process() error {
	// TODO: setup the nanoagent services
	mist := models.Service{}
	data.Get(util.AppName(), "mist", &mist)

	client, err := clients.New(mist.ExternalIP+":1445", "123")
	if err != nil {
		return err
	}
	
	if err := client.Subscribe([]string{"log"}); err != nil {
		return err
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	
	for {
		select {
		case msg := <-client.Messages():
			fmt.Printf("message: %+v\n", msg)
		case <-sigChan:
			fmt.Println("quit")
			return nil
		}
	}

	return nil
}
