package processors

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/models"
	// "github.com/nanobox-io/nanobox/util/provider"
)

// displays status about provider status and running apps
func Status() error {

	// print the header
	nameLength := longestName()
	pathLength := longestPath()

	fmt.Println()
	fmtString := fmt.Sprintf("%%-%ds : %%-7s : %%%ds\n", nameLength+10, pathLength)
	header := fmt.Sprintf(fmtString, "Status", "Running", "Path")
	fmt.Printf(header)
	fmt.Println(strings.Repeat("-", len(header)))

	envs, _ := models.AllEnvs()
	for _, env := range envs {

		apps, _ := env.Apps()
		for _, app := range apps {
			fmt.Printf(fmtString, fmt.Sprintf("%s (%s)", env.Name, app.DisplayName()), app.Status, env.Directory)
		}
	}

	// end with a newline
	fmt.Println()

	return nil
}

// returns the longest name
func longestName() int {
	longest := ""

	envs, _ := models.AllEnvs()
	for _, env := range envs {
		if len(env.Name) > len(longest) {
			longest = env.Name
		}
	}

	return len(longest)
}

// returns the longest name
func longestPath() int {
	longest := ""

	envs, _ := models.AllEnvs()
	for _, env := range envs {
		if len(env.Directory) > len(longest) {
			longest = env.Name
		}
	}

	return len(longest)
}
