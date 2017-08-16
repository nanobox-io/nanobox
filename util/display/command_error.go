package display

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/nanobox-io/nanobox-boxfile"

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
Context : %s`, parsedErr.cause, parsedErr.context)

	if parsedErr.suggest != "" {
		output = fmt.Sprintf(`%s
Suggest : %s`, output, parsedErr.suggest)
	}

	if parsedErr.output != "" {
		output = fmt.Sprintf(`%s
Output  : %s`, output, parsedErr.output)
	}

	output = fmt.Sprintf("%s\n", output)

	var appID, appName string
	env, err := models.FindEnvByID(config.EnvID())
	if err == nil {
		remote, ok := env.Remotes["default"]
		if ok {
			appID = remote.ID
			appName = remote.Name
		}
	}

	// todo: ensure this matters (seen a failed build hook have a boxfile)
	// get the raw boxfile if a processed one doesn't exist
	boxfileString := env.UserBoxfile
	if boxfileString == "" {
		// config.Boxfile() is essentially `./boxfile.yml`
		boxfileString = boxfile.NewFromPath(config.Boxfile()).String()
	}

	conf, _ := models.LoadConfig()

	// submit error to nanobox
	odin.SubmitEvent(
		"desktop#error",
		"an error occurred running nanobox desktop",
		appID,
		map[string]interface{}{
			"app-id":          appID,
			"app-name":        appName,
			"boxfile":         boxfileString,
			"context":         parsedErr.context,
			"error":           parsedErr.cause,
			"mount-type":      conf.MountType,
			"nanobox-version": models.VersionString(),
			"output":          parsedErr.output,
			"os":              runtime.GOOS,
			"provider":        conf.Provider,
			"suggest":         parsedErr.suggest,
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
	output  string
	suggest string
	team    string
}

// parseTeam parses out an error code (in this iteration, team name), we'll have added to the beginning of the "cause"
func parseTeam(err string) (team, cause string) {
	re := regexp.MustCompile(`\[([A-Z0-9]+)\] `) // matches to split after "[TEAM] "
	match := re.FindStringSubmatch(err)
	remaining := re.Split(err, 2)
	if len(match) > 0 && len(remaining) > 0 {
		return match[len(match)-1], remaining[len(remaining)-1]
	}
	return "", err
}

// parseCommandErr retrieves the cause, context, and team responsible for the error
func parseCommandErr(err error) errBits {
	var team, cause, context string

	// if it is one of our utility errors we can
	// extract the cause and the stack seperately
	if er, ok := err.(util.Err); ok {
		bits := errBits{}
		bits.team, bits.cause = parseTeam(er.Message)
		if er.Code != "" {
			bits.team = er.Code
		}
		bits.context = strings.Join(er.Stack, " -> ")
		bits.suggest = er.Suggest
		bits.output = strings.TrimSpace(er.Output)
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
