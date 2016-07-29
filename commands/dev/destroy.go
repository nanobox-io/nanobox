package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

// DestroyCmd ...
var (
	DestroyCmd = &cobra.Command{
		Use:    "destroy",
		Short:  "Destroys the docker machines associated with your dev app.",
		Long:   ``,
		PreRun: validate.Requires("provider", "provider_up"),
		Run:    destroyFn,
	}

	destroyCmdFlags = struct {
		app string
	}{}
)

//
func init() {
	DestroyCmd.Flags().StringVarP(&destroyCmdFlags.app, "app", "a", "", "app to destroy")
}

// destroyFn ...
func destroyFn(ccmd *cobra.Command, args []string) {
	appID := getAppID()
	if appID != "" {
		processor.DefaultControl.Meta["app_name"] = appID
	}
	print.OutputCommandErr(processor.Run("dev_destroy", processor.DefaultControl))
}

// look up the real app id based on what they told me.
func getAppID() string {
	if destroyCmdFlags.app == "" {
		return ""
	}
	keys, _ := data.Keys("apps")
	for _, appID := range keys {
		app := models.App{}
		data.Get("apps", appID, &app)
		if app.ID == destroyCmdFlags.app || app.Name == destroyCmdFlags.app {
			return app.ID
		}
	}

	return ""
}
