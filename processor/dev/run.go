package dev

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
)

// Run ...
type Run struct {
	App models.App

	boxfile   boxfile.Boxfile
	starts    map[string]string
	container string
}

//
func (run *Run) Run() error {
	// get the boxfile
	if err := run.loadBoxfile(); err != nil {
		return err
	}

	lumber.Debug("devRun:boxfile: %+v", run.boxfile)

	// load the start commands from the boxfile
	if err := run.loadStarts(); err != nil {
		return err
	}

	// get the id of the container we will be running in
	id := fmt.Sprintf("nanobox_%s", run.App.ID)
	if container, err := docker.GetContainer(id); err == nil {
		run.container = container.ID
	}

	// run the start commands in from the boxfile
	// in the dev container
	if err := run.runStarts(); err != nil {
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

func (run *Run) loadBoxfile() error {
	run.boxfile = boxfile.New([]byte(run.App.DeployedBoxfile))

	if !run.boxfile.Valid {
		return fmt.Errorf("the boxfile from the build is invalid")
	}
	return nil
}

func (run *Run) loadStarts() error {
	run.starts = map[string]string{}

	// loop through the nodes and get there start commands
	for _, node := range run.boxfile.Nodes("code") {

		values := run.boxfile.Node(node).Value("start")

		switch values.(type) {
		case string:
			run.starts[node] = values.(string)
		case []interface{}:
			// if it is an array we need the keys to be
			// web.site.2 where 2 is the index of the element
			for index, iFace := range values.([]interface{}) {
				if str, ok := iFace.(string); ok {
					run.starts[fmt.Sprintf("%s.%d", node, index)] = str
				}
			}
		case map[interface{}]interface{}:
			for key, val := range values.(map[interface{}]interface{}) {
				k, keyOk := key.(string)
				v, valOk := val.(string)
				if keyOk && valOk {
					run.starts[fmt.Sprintf("%s.%s", node, k)] = v
				}
			}
		}
	}
	return nil
}

func (run Run) runStarts() error {
	// loop through the starts and run them in go routines
	for key, start := range run.starts {
		go run.runStart(key, start)
	}
	return nil
}

func (run Run) runStart(name, command string) error {

	// create the docker command
	cmd := []string{
		"docker",
		"exec",
		"-u",
		"gonano",
		run.container,
		"/bin/bash",
		"-lc",
		fmt.Sprintf("cd /app/; %s", command),
	}

	lumber.Debug("run:runstarts: %+v", cmd)
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
