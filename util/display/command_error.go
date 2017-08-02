package display

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/odin"
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

	parsedErr := parseCommandErr(err)

	output := fmt.Sprintf(`
Error   : %s
Context : %s
`, parsedErr.cause, parsedErr.context)

	app := ""
	env, err := models.FindEnvByID(config.EnvID())
	if err == nil {
		remote, ok := env.Remotes["default"]
		if ok {
			app = remote.ID
		}
	}

	conf, _ := models.LoadConfig()

	// submit error to nanobox
	odin.SubmitEvent(
		"desktop#error",
		"an error occurred running nanobox desktop",
		app,
		map[string]interface{}{
			"boxfile":         env.UserBoxfile, // todo: this doesn't seem to populate
			"context":         parsedErr.context,
			"error":           parsedErr.cause,
			"mount-type":      conf.MountType,
			"nanobox-version": models.VersionString(),
			"os":              runtime.GOOS,
			"provider":        conf.Provider,
			"team":            parsedErr.team,
		},
	)

	// display error to user
	fmt.Println(output)

	if runtime.GOOS == "windows" {
		// The update process was spawned in a separate window, which will
		// close as soon as this command is finished. To ensure they see the
		// message, we need to hold open the process until they hit enter.
		fmt.Println("Enter to continue:")
		var input string
		fmt.Scanln(&input)
	}
	if exitCode == 0 {
		// todo: handle
		os.Exit(1)
	}
	// todo: handle
	os.Exit(exitCode)
}

// errBits holds the error bits we'll use for context during error reports
type errBits struct {
	cause   string
	context string
	team    string
}

// parseTeam parses out an error code (in this iteration, team name), we'll have added to the beginning of the "cause"
func parseTeam(err string) (team, cause string) {
	re := regexp.MustCompile(`\[([A-Z]+)\] `) // matches to split after "[TEAM] "
	match := re.FindStringSubmatch(err)
	remaining := re.Split(err, 2)
	return match[len(match)-1], remaining[len(remaining)-1]
}

// parseCommandErr retrieves the cause, context, and team responsible for the error
func parseCommandErr(err error) errBits {
	var team, cause, context string

	// if it is one of our utility errors we can
	// extract the cause and the stack seperately
	if er, ok := err.(util.Err); ok {
		bits := errBits{}
		bits.team, bits.cause = parseTeam(er.Message)
		bits.context = strings.Join(er.Stack, " -> ")
		return bits
	}

	trace := err.Error()
	// remove any extra : at the end of the trace
	trace = CmdErrRegex.ReplaceAllString(trace, "")

	// split the trace into a slice
	stack := strings.Split(trace, ": ")

	// extract the last item off of the list
	cause = stack[len(stack)-1]

	team, cause = parseTeam(cause)

	// now remove the last item from the list
	stack = stack[:len(stack)-1]

	// join the stack to create a context
	context = strings.Join(stack, " -> ")

	return errBits{
		cause:   cause,
		context: context,
		team:    team,
	}
}
