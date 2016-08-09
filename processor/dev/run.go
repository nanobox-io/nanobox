package dev

import (
	"os"
	"os/exec"
	"os/signal"
	"fmt"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
)

// processDevRun ...
type processDevRun struct {
	control   processor.ProcessControl
	boxfile   boxfile.Boxfile
	starts    map[string][]string
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
	devRun.starts = map[string][]string{}

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
		startSlice := devRun.boxfile.Node(node).StringSliceValue("start")
		devRun.starts[node] = startSlice
	}
	return nil
}

func (devRun processDevRun) runStarts() error {
	// loop through the starts and run them in go routines
	for key, starts := range devRun.starts {
		for _, start := range starts {
			go devRun.runStart(key, start)
		}
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
		"-it",
		devRun.container,
		"/bin/bash",
		"-c",
		fmt.Sprintf("\"%s\"", command),
	}

	fmt.Println(cmd)
	process := exec.Command(cmd[0], cmd[1:]...)

	// TODO: these will be replaced with something from the
	// new print library
	// we will also want to use 'name' to create some prefix
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	// run the process
	if err := process.Run(); err != nil && err.Error() != "exit status 137" {
		return err
	}

	return nil	
}



