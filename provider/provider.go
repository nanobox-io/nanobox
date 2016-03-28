package provider

import (
	"errors"

	"github.com/nanobox-io/nanobox/util/nanofile"
)

type (
	Provider interface {
		Info() error
		Display() error
		Reload() error
		Stop() error
		Start() error
		Init() error
	}
)

var providers = map[string]Provider{}
var invalidProvider = errors.New("invalid provider")

func Register(name string, p Provider) {
	providers[name] = p
}

func Info() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Info()
}
func Display() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Display()
}
func Reload() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Reload()
}
func Stop() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Stop()
}
func Start() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Start()
}
func Init() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Init()
}