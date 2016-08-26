package boxfile

import (
	"fmt"

	driver "github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
)

// BuildImage fetches the build image from the boxfile
func BuildImage() string {
	// first let's see if the user has a custom build image they want to use
	box := driver.NewFromPath(config.Boxfile())
	image := box.Node("build").StringValue("image")

	// then let's set the default if the user hasn't specified
	if image == "" {
		image = "nanobox/build:v1"
	}

	return image
}

// ComponentImage returns the image for the component
func ComponentImage(component *models.Component) (string, error) {
	// fetch the env
	env, err := models.FindEnvByID(component.EnvID)
	if err != nil {
		return "", fmt.Errorf("failed to load env model: %s", err.Error())
	}

	box := driver.New([]byte(env.BuiltBoxfile))
	image := box.Node(component.Name).StringValue("image")

	// the only way image can be empty is if it's a platform service
	if image == "" {
		image = fmt.Sprintf("nanobox/%s", component.Name)
	}

	return image, nil
}

// ComponentConfig returns the config data from the component boxfile
func ComponentConfig(component *models.Component) (config map[string]interface{}, err error) {

	// fetch the env
	env, err := models.FindEnvByID(component.EnvID)
	if err != nil {
		err = fmt.Errorf("failed to load env model: %s", err.Error())
		return
	}

	box := driver.New([]byte(env.BuiltBoxfile))
	config = box.Node(component.Name).Node("config").Parsed

	switch component.Name {
	case "portal", "logvac", "hoarder", "mist":
		config["token"] = "123"
	}
	return
}
