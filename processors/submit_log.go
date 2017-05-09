package processors

import (
	"fmt"
	"strings"
	"runtime"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/odin"
)

func SubmitLog(args string) error {
	// if we are running as privilage we dont submit
	if util.IsPrivileged() {
		return nil
	}

	auth, _ := models.LoadAuth()
	conf, _ := models.LoadConfig()
	if conf.CIMode {
		return nil
	}
	
	if auth.Key == "" && !conf.Anonymous {
		err :=  Login("", "", "")
		if err != nil {
			return err
		}
	}

	app := ""

	env, err := models.FindEnvByID(config.EnvID())
	if strings.Contains(args, "deploy") || strings.Contains(args, "tunnel") || strings.Contains(args, "console") {
		if err == nil {
			remote, ok := env.Remotes["default"]
			if ok {
				app = remote.ID
			}
		}	
	}

	// tell nanobox 
	go odin.SubmitEvent(
		fmt.Sprintf("desktop%s", strings.Replace(args, " ", "/", -1)),
		fmt.Sprintf("desktop command: nanobox %s", args),
		app,
		map[string]interface{}{
			"os":         runtime.GOOS,
			"provider":   conf.Provider,
			"mount-type": conf.MountType,
			"boxfile": env.UserBoxfile,
		},
	)

	return nil
}
