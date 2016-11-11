package processors

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/provider"
)

type status struct {
	envName   string
	appName   string
	status    string
	directory string
}


// displays status about provider status and running apps
func Status() error {
	fmt.Printf("Status: %s\n", provider.Status())
	fmt.Println()


	statuses := []status{}

	envs, _ := models.AllEnvs()
	for _, env := range envs {

		apps, _ := env.Apps()
		for _, app := range apps {
			statuses = append(statuses, status{
				envName:   env.Name,
				appName:   app.DisplayName(),
				status:    app.Status,
				directory: env.Directory,
				})
		}
	}

	if len(statuses) == 0 {
		return nil
	}


	// print the header
	nameLength := longestName(statuses)
	pathLength := longestPath(statuses)

	fmtString := fmt.Sprintf("%%-%ds : %%-7s : %%-%ds\n", nameLength, pathLength)

	fmt.Printf(fmtString, "App", "Status", "Path")
	fmt.Println(strings.Repeat("-", nameLength+pathLength+13))

	for _, status := range statuses {
		fmt.Printf(fmtString, fmt.Sprintf("%s (%s)", status.envName,  status.appName), status.status,  status.directory)
	}

	// end with a newline
	fmt.Println()

	return nil
}

// returns the longest name
func longestName(statuses []status) (rtn int) {

	for _, status := range statuses {
		name := fmt.Sprintf("%s (%s)", status.envName,  status.appName)
		if len(name)> rtn {
			rtn = len(name)
		}
	}

	return
}

// returns the longest name
func longestPath(statuses []status) (rtn int) {

	for _, status := range statuses {
		if len(status.directory)> rtn {
			rtn = len(status.directory)
		}
	}

	return
}
