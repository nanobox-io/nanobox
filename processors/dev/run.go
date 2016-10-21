package dev

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
)

// Run ...
func Run(appModel *models.App) error {

	display.InfoDevRunContainer(appModel.GlobalIPs["env"])

	// load the start commands from the boxfile
	starts := loadStarts(appModel)

	if len(starts) == 0 {
		display.DevRunEmpty()
		return nil
	}

	console := models.Console{"run", generator.DevName()}
	console.Save()
	defer console.Delete()

	fmt.Println()

	// run the start commands in from the boxfile
	// in the dev container
	if err := runStarts(starts); err != nil {
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

func loadStarts(appModel *models.App) map[string]string {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))
	starts := map[string]string{}
	// loop through the nodes and get there start commands
	for _, node := range boxfile.Nodes("code") {

		// print the node header
		fmt.Println(node)

		values := boxfile.Node(node).Value("start")

		switch values.(type) {
		case string:
			str := values.(string)
			starts[node] = str
			fmt.Printf("  $ %s\n", str)
		case []interface{}:
			// if it is an array we need the keys to be
			// web.site.2 where 2 is the index of the element
			for index, iFace := range values.([]interface{}) {
				if str, ok := iFace.(string); ok {
					starts[fmt.Sprintf("%s.%d", node, index)] = str
					fmt.Printf("  $ %s\n", str)
				}
			}
		case map[interface{}]interface{}:
			for key, val := range values.(map[interface{}]interface{}) {
				k, keyOk := key.(string)
				v, valOk := val.(string)
				if keyOk && valOk {
					starts[fmt.Sprintf("%s.%s", node, k)] = v
					fmt.Printf("  $ %s\n", v)
				}
			}
		case map[string]interface{}:
			for key, val := range values.(map[string]interface{}) {
				v, valOk := val.(string)
				if valOk {
					starts[fmt.Sprintf("%s.%s", node, key)] = v
					fmt.Printf("  $ %s\n", v)
				}
			}
		}
	}
	return starts
}

func runStarts(starts map[string]string) error {
	// loop through the starts and run them in go routines
	for key, start := range starts {
		go runStart(key, start)
	}
	return nil
}

func runStart(name, command string) error {

	// create the docker command
	cmd := []string{
		"-lc",
		fmt.Sprintf("cd /app/; exec %s", command),
	}

	lumber.Debug("run:runstarts: %+v", cmd)

	streamer := display.NewPrefixedStreamer("info", fmt.Sprintf("[%s] ", name))
	output, err := util.DockerExec(generator.DevName(), "gonano", "/bin/bash", cmd, streamer)
	if err != nil {
		lumber.Error("dev:runStart:util.DockerExec(%s, %s, %s, %s): %s", generator.DevName(), "gonano", "/bin/bash", cmd, err)
		return fmt.Errorf("runstart error: %s, %s", output, err)
	}

	return nil
}
