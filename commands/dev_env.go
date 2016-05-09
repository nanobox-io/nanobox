//
package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

var (
	//
	DevEnvCmd = &cobra.Command{
		Use:   "env",
		Short: "run an env command",
		Long:  ``,
	}

	DevEnvAddCmd = &cobra.Command{
		Use:   "add",
		Short: "run an env command",
		Long:  ``,
		Run: func(ccmd *cobra.Command, args []string) {
			evars := models.EnvVars{}
			data.Get(util.AppName()+"_meta", "env", &evars)
			for _, arg := range args {

				for _, pair := range strings.Split(arg, ",") {
					parts := strings.Split(pair, ":")
					if len(parts) == 2 {
						evars[strings.ToUpper(parts[0])] = parts[1]
					}
				}
			}

			data.Put(util.AppName()+"_meta", "env", evars)
		},
	}

	DevEnvRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "run an env command",
		Long:  ``,
		Run: func(ccmd *cobra.Command, args []string) {
			evars := models.EnvVars{}
			data.Get(util.AppName()+"_meta", "env", &evars)
			for _, arg := range args {
				for _, key := range strings.Split(arg, ",") {
					delete(evars, strings.ToUpper(key))
				}
			}
			data.Put(util.AppName()+"_meta", "env", evars)
		},
	}

	DevEnvListCmd = &cobra.Command{
		Use:   "list",
		Short: "run an env command",
		Long:  ``,
		Run: func(ccmd *cobra.Command, args []string) {
			evars := models.EnvVars{}
			data.Get(util.AppName()+"_meta", "env", &evars)
			fmt.Println(evars)

		},
	}
)

func init() {
	DevEnvCmd.AddCommand(DevEnvAddCmd)
	DevEnvCmd.AddCommand(DevEnvRemoveCmd)
	DevEnvCmd.AddCommand(DevEnvListCmd)
}
