package link

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/odin"
)

func Add(envModel *models.Env, appName, alias string) error {

	// ensure the env model has been generated
	if err := envModel.Generate(); err != nil {
		return fmt.Errorf("unable to generate the environment")
	}

	// set the alias to be the default its missing
	if alias == "" {
		alias = "default"
	}

	// set the appName to the folder name if its missing
	if appName == "" {
		appName = config.LocalDirName()
	}

	// fetch the odin app
	app, err := odin.App(appName)
	if err != nil {
		fmt.Printf("! Sorry, but you don't have access to %s\n", appName)
		return nil
	}

	// ensure the links map is initialized
	if envModel.Links == nil {
		envModel.Links = map[string]models.Link{}
	}

	envModel.Links[alias] = models.Link{app.ID, app.Name}

	if err := envModel.Save(); err != nil {
		return fmt.Errorf("failed to save link: %s", err.Error())
	}

	fmt.Printf("%s Codebase linked to %s\n", display.TaskComplete, appName)
	
	if alias != "default" {
		fmt.Printf("  through the '%s' alias\n", alias)
	}

	return nil
}
