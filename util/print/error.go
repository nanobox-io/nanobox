package print

import (
	"fmt"
	"os"

	"github.com/pagodabox/nanobox-golang-stylish"
)

// OutputCommandErr ...
//
// We hit a minor bump, which can be quickly resolved following the instructions above.
// If not, come talk and we'll walk you through the resolution.
func OutputCommandErr(err error) {
	if err != nil {
		fmt.Printf(`

Whoops, looks like we encountered a small error:

%s

Following the instructions above should help resolve it quickly. If you're still
experiencing issues, please come talk with us and we'll work it out together.

irc   : #nanobox
email : help@nanobox.io

`, err.Error())
		os.Exit(1)
	}
}

// OutputProcessorErr ...
func OutputProcessorErr(heading, message string) {
	fmt.Printf(stylish.Error(heading, message))
}
