package processors

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/provider"
)

// displays status about provider status and running apps
func Status() error {

	// print the header
	leftWidth := len(longestName()) + 7
	margin := strings.Repeat(" ", leftWidth-5)
	fmt.Printf("\nStatus%s: %s\n", margin, provider.Status())

	margin = strings.Repeat("-", leftWidth+10)
	fmt.Printf("%s\n", margin)

	envs, _ := models.AllEnvs()
	for _, env := range envs {

		// calculate a margin
		margin := strings.Repeat(" ", leftWidth-len(env.Name)-5)

		apps, _ := env.Apps()
		for _, app := range apps {
			fmt.Printf("%s (%s)%s: %s\n", env.Name, app.Name, margin, app.Status)
		}
	}

	// end with a newline
	fmt.Println()

	return nil
}

// returns the longest name
func longestName() string {
	longest := ""

	envs, _ := models.AllEnvs()
	for _, env := range envs {
		if len(env.Name) > len(longest) {
			longest = env.Name
		}
	}

	return longest
}
