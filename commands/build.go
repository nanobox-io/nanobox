package commands

import (
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/generators/hooks/build"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (
	// BuildCmd builds the app's runtime.
	BuildCmd = &cobra.Command{
		Use:   "build-runtime",
		Short: "Build your app's runtime.",
		Long: `
Builds your app's runtime, which is used both
locally and in live environments.
		`,
		PreRun:  steps.Run("start"),
		Run:     buildFn,
		Aliases: []string{"build"},
	}

	cacheClear bool
)

func init() {
	steps.Build("build-runtime", buildComplete, buildFn)

	BuildCmd.Flags().BoolVar(&cacheClear, "clear-cache", false, "Clear package cache for this build.")
}

func buildFn(ccmd *cobra.Command, args []string) {
	if cacheClear {
		build.ClearPkgCache = true
	}

	env, _ := models.FindEnvByID(config.EnvID())
	display.CommandErr(processors.Build(env))
}

// update: this runs on deploy
func buildComplete() bool {
	// check the boxfile to be sure it hasnt changed
	env, _ := models.FindEnvByID(config.EnvID())
	box := boxfile.NewFromPath(config.Boxfile())

	// we need to rebuild if this isnt true without going to check triggers
	if env.UserBoxfile == "" || env.UserBoxfile != box.String() {
		return false
	}

	// now check to see if any of the build triggers have changed
	lastBuildsBoxfile := boxfile.New([]byte(env.BuiltBoxfile))
	for _, trigger := range lastBuildsBoxfile.Node("run.config").StringSliceValue("build_triggers") {
		if env.BuildTriggers[trigger] != util.FileMD5(trigger) {
			return false
		}
	}
	return true
}
