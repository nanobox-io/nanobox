//
package terminal

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/term"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/config"
)

//
func PrintNanoboxHeader(kind string) {
	switch kind {

	//
	case "exec":
		os.Stderr.WriteString(stylish.Bullet("Executing command in nanobox..."))

		//
	case "develop", "console":
		os.Stderr.WriteString(`+> Opening a nanobox console:


                                     **
                                  ********
                               ***************
                            *********************
                              *****************
                            ::    *********    ::
                               ::    ***    ::
                             ++   :::   :::   ++
                                ++   :::   ++
                                   ++   ++
                                      +

                      _  _ ____ _  _ ____ ___  ____ _  _
                      |\ | |__| |\ | |  | |__) |  |  \/
                      | \| |  | | \| |__| |__) |__| _/\_
`)

		if kind == "develop" {
			os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ You are in a virtual machine (vm)
+ Your local source code has been mounted into the vm, and changes in either
the vm or local will be mirrored.
+ If you run a server, access it at >> %s
--------------------------------------------------------------------------------
`, config.Nanofile.Domain))
		}
	}
}

// GetTTYSize
func GetTTYSize(fd uintptr) (int, int) {

	ws, err := term.GetWinsize(fd)
	if err != nil {
		config.Fatal("[util/server/exec] term.GetWinsize() failed", err.Error())
	}

	//
	return int(ws.Width), int(ws.Height)
}
