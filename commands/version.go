package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
)

var (

	// VersionCmd ...
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the current Nanobox version.",
		Long:  ``,
		Run:   versionFn,
	}
)

// versionFn ...
func versionFn(ccmd *cobra.Command, args []string) {
	v := "2.0.0"
	update, _ := models.LoadUpdate()
	md5Parts := strings.Fields(update.CurrentVersion)
	md5 := ""
	if len(md5Parts) > 1 {
		md5 = md5Parts[len(md5Parts)-1]
	}
	fmt.Printf("Nanobox version %s (%s)\n", v, md5)
}
