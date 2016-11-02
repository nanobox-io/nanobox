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
		Short: "version",
		Long:  ``,
		Run:   versionFn,
	}
)

// versionFn ...
func versionFn(ccmd *cobra.Command, args []string) {
	v := "0.9.0"
	update, _ := models.LoadUpdate()
	md5Parts := strings.Fields(update.CurrentVersion)
	fmt.Printf("Nanobox version %s (%s)\n", v, md5Parts[len(md5Parts)-1])
}
