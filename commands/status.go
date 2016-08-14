package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/provider"
)

var (

	// EnvCmd ...
	StatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Displays the status of your Nanobox VM & running platforms.",
		Long:  ``,
		Run:   statusFn,
	}
)

func statusFn(ccmd *cobra.Command, args []string) {
	fmt.Printf("Provider status: %s\n", provider.Status())
	envs, _ := models.AllEnvs()
	for _, env := range envs {
		fmt.Printf("Environment: %s", env.Name)
		apps, _ := models.AllAppsByEnv(env.ID)
		for _, app := range apps {
			fmt.Println("  %s(%s): %s\n", app.Name, app.ID, app.Status)
		}
	}
}
