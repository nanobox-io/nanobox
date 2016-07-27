package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/data"
)

var (

	// EnvCmd ...
	StatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Manages environment variables in your local dev app.",
		Long:  ``,
		Run: statusFn,
	}
)

func statusFn(ccmd *cobra.Command, args []string) {
	fmt.Printf("Provider status: %s\n", provider.Status())
	appNames, _ := data.Keys("apps")
	for _, appName := range appNames {
		app := models.App{}
		data.Get("apps", appName, &app)
		fmt.Printf("  %s: %s\n", app.Name, app.Status)
	}
}

