package provider

import (
	"errors"

	"github.com/nanobox-io/nanobox/util/nanofile"
)

type (
	Provider interface {
		Create() error
		Reboot() error
		Stop() error
		Destroy() error
		Start() error
		AddIP(ip string) error
		RemoveIP(ip string) error
		AddNat(host, container string) error
		RemoveNat(host, container string) error
		AddMount(local, host string) error
		RemoveMount(local, host string) error
	}
)

var providers = map[string]Provider{}
var invalidProvider = errors.New("invalid provider")

func Register(name string, p Provider) {
	providers[name] = p
}

func Create() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Create()
}
func Reboot() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Reboot()
}
func Stop() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Stop()
}
func Destroy() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Destroy()
}
func Start() error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.Start()
}
func AddIP(ip string) error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.AddIP(ip)
}
func RemoveIP(ip string) error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.RemoveIP(ip)
}
func AddNat(host, container string) error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.AddNat(host, container)
}
func RemoveNat(host, container string) error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.RemoveNat(host, container)
}
func AddMount(local, host string) error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.AddMount(local, host)
}
func RemoveMount(local, host string) error {
	p, ok := providers[nanofile.Viper().GetString("provider")]
	if !ok {
		return invalidProvider
	}
	return p.RemoveMount(local, host)
}
