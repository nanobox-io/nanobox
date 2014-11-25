package commands

import (
	"fmt"
	// "strconv"

	nanoAPI "github.com/nanobox-core/api-client-go"
	// "github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// EVarListCommand satisfies the Command interface for listing an app's environment
// variables
type EVarListCommand struct{}

// Help prints detailed help text for the evar list command
func (c *EVarListCommand) Help() {
	ui.CPrintln(`
Description:
  Lists an app's environment variables

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

  type:
    'Custom' evar's are preceded by a '*'.

Usage:
  pagoda evar:list [-a app-name]

  ex. pagoda evar:list -a app-name

Options:
  -a, --app [app-name]
    The name of the app
  `)
}

type Test struct{}

// Run disaplys select information about all of an app's environment variables
func (c *EVarListCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	// if fApp == "" {
	// 	fApp = helpers.FindPagodaApp()
	// }

	fApp = "TESTING"

	// get environment variables
	eVars, err := api.GetEVars()
	if err != nil {
		fmt.Printf("There was a problem getting '%v's' environment variables. See ~/.pagodabox/log.txt for details", fApp)
		ui.Error("pagoda evar:list", err)
	}

	fmt.Println("INITIAL DONE!!!", eVars)

	thing := Test{}

	fmt.Println("DO MIST!!!")

	//
	api.DoRawRequest(&thing, "GET", "http://127.0.0.1:1445/mist?subscribe=a,b", nil, nil)

	fmt.Println("MIST DONE!!!")

	// 	var internal, custom []nanoAPI.EVar

	// 	for _, eVar := range eVars {

	// 		// load custom environment variables
	// 		if !eVar.Internal {
	// 			custom = append(custom, eVar)

	// 			// load generated environment variables
	// 		} else {
	// 			internal = append(internal, eVar)
	// 		}
	// 	}

	// 	fmt.Println(`
	// Custom (` + strconv.Itoa(len(custom)) + `):
	// --------------------------------------------------`)

	// 	// list custom environment variables
	// 	if len(custom) > 0 {
	// 		for _, eVar := range custom {
	// 			fmt.Printf("%v = %v", eVar.Title, eVar.Value)
	// 			fmt.Println("")
	// 		}
	// 	} else {
	// 		fmt.Println("** NONE CREATED **")
	// 	}

	// 	fmt.Println("")

	// 	// list generated environment variables
	// 	fmt.Println(`
	// Generated (` + strconv.Itoa(len(internal)) + `):
	// --------------------------------------------------`)
	// 	for _, eVar := range internal {
	// 		fmt.Printf("%v = %v", eVar.Title, eVar.Value)
	// 		fmt.Println("")
	// 	}

	// 	fmt.Println("")
}
