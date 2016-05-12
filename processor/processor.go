package processor

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"	
	"github.com/nanobox-io/nanobox/util/data"	
	"github.com/nanobox-io/nanobox/models"
)

type (
	ProcessConfig struct {
		DevMode    bool
		Verbose    bool
		Background bool
		Force      bool
		Meta       map[string]string
	}

	ProcessBuilder func(ProcessConfig) (Processor, error)

	Processor interface {
		Process() error
		Results() ProcessConfig
	}
)

var (
	DefaultConfig = ProcessConfig{Meta: map[string]string{}}
	processors    = map[string]ProcessBuilder{}
)

func Register(name string, sb ProcessBuilder) {
	if _, ok := processors[name]; ok {
		panic("Duplicate Registration - " + name)
	}
	processors[name] = sb
}

func Build(name string, pc ProcessConfig) (Processor, error) {
	lumber.Debug(name)
	proc, ok := processors[name]
	if !ok {
		return nil, fmt.Errorf("Invalid Processor %s", name)
	}
	return proc(pc)
}

func Run(name string, pc ProcessConfig) error {
	proc, err := Build(name, pc)
	if err != nil {
		return err
	}
	return proc.Process()
}

func getAppID(alias string) string {
	link := models.AppLinks{}
	data.Get(util.AppName()+"_meta", "links", &link)
	if alias == "" {
		alias = "default"
	}
	app, ok := link[alias]
	if !ok {
		return alias
	}
	return app
}
