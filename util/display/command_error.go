package display

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/nanobox-io/nanobox/commands/registry"
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
	// get the exit code we are going to use
	// if none has been set GetInt returns 0
	exitCode := registry.GetInt("exit_code")

	if err == nil {
		// if an exit code is provided we need to quit here
		// and use that exit code
		if exitCode != 0 {
			os.Exit(exitCode)
		}
		return
	}

	cause, context := parseCommandErr(err)

	fmt.Println()
	fmt.Printf("Error   : %s\n", cause)
	fmt.Printf("Context : %s\n", context)
	fmt.Println()

	if runtime.GOOS == "windows" {
		// The update process was spawned in a separate window, which will
		// close as soon as this command is finished. To ensure they see the
		// message, we need to hold open the process until they hit enter.
		fmt.Println("Enter to continue:")
		var input string
		fmt.Scanln(&input)
	}
	if exitCode == 0 {
		os.Exit(1)
	}
	os.Exit(exitCode)
}

func parseCommandErr(err error) (cause, context string) {
	// if it is one of our utility errors we can
	// extract the cause and the stack seperately
	if er, ok := err.(util.Err); ok {
		return er.Message, strings.Join(er.Stack, " -> ")
	}

	trace := err.Error()
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
