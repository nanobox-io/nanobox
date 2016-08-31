package link

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/odin"
)

// Add
func Add(envModel *models.Env, appName, alias string) error {

	if err := envModel.Generate(); err != nil {
		return fmt.Errorf("unable to generate the environment")
	}

	// set the alias to be the default its missing
	if alias == "" {
		alias = "default"
	}

	// set the appName to the folder name its missing
	if appName == "" {
		appName = config.LocalDirName()
	}

	// get app id
	app, err := odin.App(appName)
	if err != nil {
		return err
	}

	if envModel.Links == nil {
		envModel.Links = map[string]string{}
	}

	envModel.Links[alias] = app.ID

	return envModel.Save()
}
