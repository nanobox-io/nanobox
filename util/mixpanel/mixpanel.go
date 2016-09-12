package mixpanel

import (
	"runtime"

	"github.com/jcelliott/lumber"
	mp "github.com/timehop/go-mixpanel"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
)

var token string

func Report(args string) {
	go func() {

		mx := mp.NewMixpanel(token)
		id := util.UniqueID()

		err := mx.Track(id, "command", mp.Properties{
			"os":         runtime.GOOS,
			"provider":   config.Viper().GetString("provider"),
			"mount-type": config.Viper().GetString("mount-type"),
			"args":       args,
			"cpus":       runtime.NumCPU(),
		})

		if err != nil {
			lumber.Error("mixpanel(%s).Report(%s): %s", token, args, err)
		}
	}()
}
