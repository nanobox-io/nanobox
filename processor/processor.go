package processor

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

type (

	// BreadCrumbProcessor ...
	BreadCrumbProcessor struct {
		crumb string
		proc  Processor
	}

	// ProcessControl ...
	ProcessControl struct {
		DevMode      bool
		Quiet        bool
		Verbose      bool
		Background   bool
		Force        bool
		DisplayLevel int
		Meta         map[string]string
	}

	// ProcessBuilder ...
	ProcessBuilder func(ProcessControl) (Processor, error)

	// Processor ...
	Processor interface {
		Process() error
		Results() ProcessControl
	}
)

var (
	// DefaultConfig ...
	DefaultConfig = ProcessControl{Meta: map[string]string{}}

	processors = map[string]ProcessBuilder{}
)

// Register ...
func Register(name string, sb ProcessBuilder) {
	if _, ok := processors[name]; !DefaultConfig.Force && ok {
		panic("Duplicate Registration - " + name)
	}

	//
	processors[name] = sb
}

// Build ...
func Build(name string, pc ProcessControl) (Processor, error) {
	lumber.Debug(name)
	procFunc, ok := processors[name]
	if !ok {
		return nil, fmt.Errorf("Invalid Processor %s", name)
	}
	proc, err := procFunc(pc)
	return BreadCrumbProcessor{name, proc}, err
}

// Run ...
func Run(name string, pc ProcessControl) error {
	proc, err := Build(name, pc)
	if err != nil {
		return err
	}
	return proc.Process()
}

// ExecWriter ...
func ExecWriter() io.Writer {
	if DefaultConfig.Quiet {
		return nil
	}
	return os.Stdout
}

// Process ...
func (bcp BreadCrumbProcessor) Process() error {
	err := bcp.proc.Process()
	if err != nil {
		err = fmt.Errorf("%s:%s", bcp.crumb, err.Error())
	}
	return err
}

// Results ...
func (bcp BreadCrumbProcessor) Results() ProcessControl {
	return bcp.proc.Results()
}

// Display displays all of the possible
func (control ProcessControl) Display(msg string) {
	fmt.Print(stylish.Nest(control.DisplayLevel, msg))
}

// Info ...
func (control ProcessControl) Info(msg string) {
	if !DefaultConfig.Quiet {
		fmt.Print(stylish.Nest(control.DisplayLevel, msg))
	}
}

// Trace ...
func (control ProcessControl) Trace(msg string) {
	if DefaultConfig.Verbose {
		fmt.Print(stylish.Nest(control.DisplayLevel, msg))
	}
}

// getAppID ...
func getAppID(alias string) string {
	link := models.AppLinks{}
	data.Get(config.AppName()+"_meta", "links", &link)
	if alias == "" {
		alias = "default"
	}
	app, ok := link[alias]
	if !ok {
		return alias
	}

	return app
}

// connect ...
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
