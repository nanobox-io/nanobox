package print

import (
	"fmt"

	"github.com/pagodabox/nanobox-golang-stylish"
)

// OutputCommandErr ...
//
// We hit a minor bump, which can be quickly resolved following the instructions above.
// If not, come talk and we'll walk you through the resolution.
func OutputCommandErr(err error) {
	if err != nil {
		fmt.Printf(`

Whoops, looks like there was a minor bump. Following the instructions above should
help resolve it quickly.

If the instructions aren't helping, please come talk with us and we can work it
out together.

irc   : #nanobox
email : help@nanobox.io

(%s)
`, err.Error())
	}
}

// OutputProcessorErr ...
func OutputProcessorErr(heading, message string) {
	fmt.Printf(stylish.Error(heading, message))
}
