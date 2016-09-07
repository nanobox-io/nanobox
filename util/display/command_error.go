package display

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/bugsnag/bugsnag-go"
	
	"github.com/nanobox-io/nanobox/util/config"
)

// CommandErr ...
//
// We hit a minor bump, which can be quickly resolved following the instructions above.
// If not, come talk and we'll walk you through the resolution.
func CommandErr(err error) {
	if err != nil {
		lumber.Error("Command: %+s", err.Error())

		bugsnagErr := bugsnag.Notify(err, bugsnag.User{Id: config.Viper().GetString("token")}, bugsnag.SeverityInfo)
		if bugsnagErr != nil {
			lumber.Error("Bugsnag error: %s", bugsnagErr)
		}
		
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
