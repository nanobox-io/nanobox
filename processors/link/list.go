package link

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/models"
)

// List ...
func List(env *models.Env) error {

	if len(env.Links) == 0 {
		fmt.Printf("\n! This codebase is not linked to any apps\n\n")
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
	for alias, link := range env.Links {
		margin := strings.Repeat(" ", leftColWidth-len(link.Name))
		fmt.Printf("%s%s: %s\n", link.Name, margin, alias)
	}

	// end with a newline
	fmt.Println()

	return nil
}

// returns the longest name
func longestName(env *models.Env) string {
	longest := ""

	for _, link := range env.Links {
		if len(link.Name) > len(longest) {
			longest = link.Name
		}
	}

	return longest
}

// returns the longest alias
func longestAlias(env *models.Env) string {
	longest := ""

	for alias, _ := range env.Links {
		if len(alias) > len(longest) {
			longest = alias
		}
	}

	return longest
}
