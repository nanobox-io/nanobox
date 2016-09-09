package mixpanel

import (
	"runtime"
	
	mp "github.com/timehop/go-mixpanel"
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/config"	
)

var token string

func Report(args string) {
	go func() {

		mx := mp.NewMixpanel(token)
		id := config.Viper().GetString("token")

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