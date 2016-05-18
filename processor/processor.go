package processor

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
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
	_, ok := processors[name]
	if !DefaultConfig.Force && ok {
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

func connect(req *http.Request) (net.Conn, []byte, error) {

	//
	b := make([]byte, 1)

	// if we can't connect to the server, lets bail out early
	conn, err := tls.Dial("tcp4", location, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return conn, b, err
	}
	// we dont defer a conn.Close() here because we're returning the conn and
	// want it to remain open

	// make an http request

	if err := req.Write(conn); err != nil {
		return conn, b, err
	}

	// wait trying to read from the connection until a single read happens (blocking)
	if _, err := conn.Read(b); err != nil {
		return conn, b, err
	}

	return conn, b, nil
}
