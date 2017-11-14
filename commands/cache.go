package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/cache"
	"github.com/nanobox-io/nanobox/util/display"
)

func init() {
	CacheCmd.AddCommand(clearCmd)
}

var (
	// CacheCmd provides interaction with the package cache.
	CacheCmd = &cobra.Command{
		Use:   "cache",
		Short: "Manage the package cache.",
	}

	// clearCmd allows clearing the cache
	clearCmd = &cobra.Command{
		Use:   "clear [app-name]",
		Short: "Clear the package cache.",
		Long: `
The 'app-name' comes from the directory name, or 'nanobox status'.
Not specifying an app name will clear every app's pkg cache.
`,
		PreRun: steps.Run("start"),
		Run:    clearFn,
	}
)

func clearFn(ccmd *cobra.Command, args []string) {
	if len(args) > 0 {
		env, err := models.FindEnvByName(args[0])
		if env == nil {
			fmt.Printf("Failed to search apps: %s\n", err.Error())
			return
		}
		if err == nil {
			err = cache.Clear(env.ID)
			if err != nil {
				display.CommandErr(err)
				return
			}

			fmt.Printf("Cleared package cache for '%s'.\n", env.ID)
			return
		}
	}

	envs, err := models.AllEnvs()
	if err != nil {
		fmt.Printf("Failed listing apps: %s\n", err.Error())
		return
	}

	err = cache.ClearAll(envs)
	if err != nil {
		display.CommandErr(err)
		return
	}
	fmt.Println("Cleared all app package caches.")
}
