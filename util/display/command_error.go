package display

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bugsnag/bugsnag-go"
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
)

var (
	CmdErrRegex = regexp.MustCompile(":\\s?$")
)

// CommandErr ...
//
// We hit a minor bump, which can be quickly resolved following the instructions above.
// If not, come talk and we'll walk you through the resolution.
func CommandErr(err error) {

	if err == nil {
		return
	}

	lumber.Error("Command: %+s", err.Error())

	bugsnagErr := bugsnag.Notify(err, bugsnag.User{Id: util.UniqueID()}, bugsnag.SeverityInfo)
	if bugsnagErr != nil {
		lumber.Error("Bugsnag error: %s", bugsnagErr)
	}

	cause, context := parseCommandErr(err.Error())

	fmt.Println()
	fmt.Printf("Error   : %s\n", cause)
	fmt.Printf("Context : %s\n", context)
	fmt.Println()

	os.Exit(1)
}

func parseCommandErr(trace string) (cause, context string) {
	// remove any extra : at the end of the trace
	trace = CmdErrRegex.ReplaceAllString(trace, "")

	// split the trace into a slice
	stack := strings.Split(trace, ": ")

	// extract the last item off of the list
	cause = stack[len(stack)-1]

	// now remove the last item from the list
	stack = stack[:len(stack)-1]

	// join the stack to create a context
	context = strings.Join(stack, " -> ")

	return
}
