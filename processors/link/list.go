package link

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/models"
)

// List ...
func List(env *models.Env) error {

	// set the left column width to the longest name
	leftColWidth := len(longestName(env)) + 2
	
	// unless the longest name is less than 10 characters :)
	if leftColWidth < 10 {
		leftColWidth = 10
	}

	// print the header
	margin := strings.Repeat(" ", leftColWidth - 8)
	fmt.Printf("\nApp Name%s: Alias\n", margin)
	separater := strings.Repeat("-", leftColWidth + len(longestAlias(env)) + 1)
	fmt.Printf("%s\n", separater)
	
	// print the table
	for name, alias := range env.Links {
		margin := strings.Repeat(" ", leftColWidth - len(name))
		fmt.Printf("%s%s: %s\n", name, margin, alias)
	}

	// end with a newline
	fmt.Println()

	return nil
}

// returns the longest name
func longestName(env *models.Env) string {
  longest := ""
  
  for name, _ := range env.Links {
    if len(name) > len(longest) {
      longest = name
    }
  }
  
  return longest
}

// returns the longest alias
func longestAlias(env *models.Env) string {
	longest := ""
	
	for _, alias := range env.Links {
		if len(alias) > len(longest) {
			longest = alias
		}
	}
	
	return longest
}
