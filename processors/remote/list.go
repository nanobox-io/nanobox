package remote

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/models"
)

// List ...
func List(env *models.Env) error {

	if len(env.Remotes) == 0 {
		fmt.Printf("\n! This codebase is not connected to any apps\n\n")
		return nil
	}

	// set the left column width to the longest name
	leftColWidth := len(longestName(env)) + 2

	// unless the longest name is less than 10 characters :)
	if leftColWidth < 10 {
		leftColWidth = 10
	}

	// print the header
	margin := strings.Repeat(" ", leftColWidth-8)
	fmt.Printf("\nApp Name%s: Alias\n", margin)
	separator := strings.Repeat("-", leftColWidth+len(longestAlias(env))+1)
	fmt.Printf("%s\n", separator)

	// print the table
	for alias, remote := range env.Remotes {
		margin := strings.Repeat(" ", leftColWidth-len(remote.Name))
		fmt.Printf("%s%s: %s\n", remote.Name, margin, alias)
	}

	// end with a newline
	fmt.Println()

	return nil
}

// returns the longest name
func longestName(env *models.Env) string {
	longest := ""

	for _, remote := range env.Remotes {
		if len(remote.Name) > len(longest) {
			longest = remote.Name
		}
	}

	return longest
}

// returns the longest alias
func longestAlias(env *models.Env) string {
	longest := ""

	for alias := range env.Remotes {
		if len(alias) > len(longest) {
			longest = alias
		}
	}

	return longest
}
