package provider

import (
	"errors"

	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/validate"
)

type (
	Provider interface {
		HostShareDir() string
		HostMntDir() string
		Valid() error
		Create() error
		Reboot() error
		Stop() error
		Destroy() error
		Start() error
		DockerEnv() error
		AddIP(ip string) error
		RemoveIP(ip string) error
		AddNat(host, container string) error
		RemoveNat(host, container string) error
		AddShare(local, host string) error
		RemoveShare(local, host string) error
		AddMount(local, host string) error
		RemoveMount(local, host string) error
	}
)

var providers = map[string]Provider{}
var verbose = true

func Register(name string, p Provider) {
	providers[name] = p
}

func init() {
	validate.Register("provider", Valid)
}

func Display(verb bool) {
	verbose = verb
}

func Valid() error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Valid()
}

func Create() error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Create()
}

func Reboot() error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Reboot()
}

func Stop() error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Stop()
}

func Destroy() error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Destroy()
}

func Start() error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Start()
}

func HostShareDir() string {
	p, err := fetchProvider()
	if err != nil {
		return ""
	}
	return p.HostShareDir()
}

func HostMntDir() string {
	p, err := fetchProvider()
	if err != nil {
		return ""
	}
	return p.HostMntDir()
}

func DockerEnv() error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.DockerEnv()
}

func AddIP(ip string) error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddIP(ip)
}

func RemoveIP(ip string) error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveIP(ip)
}

func AddNat(host, container string) error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddNat(host, container)
}

func RemoveNat(host, container string) error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveNat(host, container)
}

func AddShare(local, host string) error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddShare(local, host)
}

func RemoveShare(local, host string) error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveShare(local, host)
}

func AddMount(local, host string) error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddMount(local, host)
}

func RemoveMount(local, host string) error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveMount(local, host)
}

// fetchProvider fetches the registered provider from the configured name
func fetchProvider() (Provider, error) {
	p, ok := providers[config.Viper().GetString("provider")]
	if !ok {
		return nil, errors.New("invalid provider")
	}

	return p, nil
}
