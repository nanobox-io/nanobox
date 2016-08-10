package dev

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processDevRun ...
type processDevRun struct {
	control   processor.ProcessControl
	boxfile   boxfile.Boxfile
	starts    map[string]string
	container string
}

//
func init() {
	processor.Register("dev_run", devRunFn)
}

//
func devRunFn(control processor.ProcessControl) (processor.Processor, error) {
	devRun := &processDevRun{control: control}
	return devRun, devRun.validateMeta()
}

//
func (devRun processDevRun) Results() processor.ProcessControl {
	return devRun.control
}

//
func (devRun processDevRun) Process() error {
	// get the boxfile
	if err := devRun.loadBoxfile(); err != nil {
		return err
	}

	// load the start commands from the boxfile
	if err := devRun.loadStarts(); err != nil {
		return err
	}

	// get the id of the container we will be running in
	id := fmt.Sprintf("nanobox_%s_dev", config.AppID())
	if container, err := docker.GetContainer(id); err == nil {
		devRun.container = container.ID
	}

	// run the start commands in from the boxfile
	// in the dev container
	if err := devRun.runStarts(); err != nil {
		return err
	}

	// catch signals and stop the run commands on signal
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt)
	defer signal.Stop(sigs)

	for range sigs {
		// if we get a interupt we will jut return here
		// causing the container to be destroyed and our
		// exec processes to die
		return nil
	}

	return nil
}

func (devRun *processDevRun) validateMeta() error {
	devRun.starts = map[string]string{}

	// currently no error conditions exist for this processor
	return nil
}

func (devRun *processDevRun) loadBoxfile() error {
	// get the build boxfile from the database
	boxfileModel := models.Boxfile{}
	if err := data.Get(config.AppID()+"_meta", "build_boxfile", &boxfileModel); err != nil {
		return fmt.Errorf("No build has been completed for this application")
	}

	// load the boxfile into the boxfile package and make sure its
	// valid
	devRun.boxfile = boxfile.New(boxfileModel.Data)
	if !devRun.boxfile.Valid {
		return fmt.Errorf("the boxfile from the build is invalid")
	}
	return nil
}

func (devRun processDevRun) loadStarts() error {
	// loop through the nodes and get there start commands
	for _, node := range devRun.boxfile.Nodes("code") {

		values := devRun.boxfile.Node(node).Value("start")

		switch values.(type) {
		case string:
			devRun.starts[node] = values.(string)
		case []interface{}:
			// if it is an array we need the keys to be
			// web.site.2 where 2 is the index of the element
			for index, iFace := range values.([]interface{}) {
				if str, ok := iFace.(string); ok {
					devRun.starts[fmt.Sprintf("%s.%d", node, index)] = str
				}
			}
		case map[interface{}]interface{}:
			for key, val := range values.(map[interface{}]interface{}) {
				k, keyOk := key.(string)
				v, valOk := val.(string)
				if keyOk && valOk {
					devRun.starts[fmt.Sprintf("%s.%s", node, k)] = v
				}
			}
		}
	}
	return nil
}

func (devRun processDevRun) runStarts() error {
	// loop through the starts and run them in go routines
	for key, start := range devRun.starts {
		go devRun.runStart(key, start)
	}
	return nil
}

func (devRun processDevRun) runStart(name, command string) error {

	// create the docker command
	cmd := []string{
		"docker",
		"exec",
		"-u",
		"gonano",
		devRun.container,
		"/bin/bash",
		"-lc",
		fmt.Sprintf("cd /app/; %s", command),
	}

	process := exec.Command(cmd[0], cmd[1:]...)

	// TODO: these will be replaced with something from the
	// new print library
	// we will also want to use 'name' to create some prefix
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	// run the process
	if err := process.Run(); err != nil && err.Error() != "exit status 137" {
		return err
	}

	return nil
}
