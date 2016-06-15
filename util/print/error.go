package print

import "fmt"

// OutputCmdErr ...
//
// We hit a minor bump, which can be quickly resolved following the instructions above.
// If not, come talk and we'll walk you through the resolution.
func OutputCmdErr(err error) {
	if err != nil {
		fmt.Printf(`

We hit a minor bump, which should be quickly resolved following the instructions
above.

If the instructions aren't helping, please come talk with us and we can work it
out together.

irc   : #nanobox
email : help@nanobox.io

(%s)
`, err.Error())
	}
}
