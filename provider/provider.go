package provider

import (
	"errors"

	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/validate"
)

// Provider ...
type Provider interface {
	HostShareDir() string
	HostMntDir() string
	HostIP() (string, error)
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
	HasShare(local, host string) bool
	AddShare(local, host string) error
	RemoveShare(local, host string) error
	HasMount(mount string) bool
	AddMount(local, host string) error
	RemoveMount(local, host string) error
	Run(command []string) ([]byte, error)
}

var (
	providers = map[string]Provider{}
	verbose   = true
)

// Register ...
func Register(name string, p Provider) {
	providers[name] = p
}

// init ...
func init() {
	validate.Register("provider", Valid)
}

// Display ...
func Display(verb bool) {
	verbose = verb
}

// Valid ...
func Valid() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Valid()
}

// Create ...
func Create() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Create()
}

// Reboot ...
func Reboot() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Reboot()
}

// Stop ...
func Stop() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Stop()
}

// Destroy ..
func Destroy() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Destroy()
}

// Start ..
func Start() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Start()
}

// HostShareDir ...
func HostShareDir() string {

	p, err := fetchProvider()
	if err != nil {
		return ""
	}

	return p.HostShareDir()
}

// HostMntDir ..
func HostMntDir() string {

	p, err := fetchProvider()
	if err != nil {
		return ""
	}

	return p.HostMntDir()
}

// HostIP ..
func HostIP() (string, error) {

	p, err := fetchProvider()
	if err != nil {
		return "", err
	}

	return p.HostIP()
}

// DockerEnv ..
func DockerEnv() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.DockerEnv()
}

// AddIP ..
func AddIP(ip string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddIP(ip)
}

// RemoveIP ...
func RemoveIP(ip string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveIP(ip)
}

// AddNat ..
func AddNat(host, container string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddNat(host, container)
}

// RemoveNat ..
func RemoveNat(host, container string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveNat(host, container)
}

func HasShare(local, host string) bool {

	p, err := fetchProvider()
	if err != nil {
		return false
	}

	return p.HasShare(local, host)
}

// AddShare ...
func AddShare(local, host string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddShare(local, host)
}

// RemoveShare ...
func RemoveShare(local, host string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveShare(local, host)
}

// HasMount ...
func HasMount(path string) bool {

	p, err := fetchProvider()
	if err != nil {
		return false
	}

	return p.HasMount(path)
}

// AddMount ...
func AddMount(local, host string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddMount(local, host)
}

// RemoveMount ...
func RemoveMount(local, host string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveMount(local, host)
}

// Run a command inside of the provider context
func Run(command []string) ([]byte, error) {

	p, err := fetchProvider()
	if err != nil {
		return nil, err
	}

	return p.Run(command)
}

// fetchProvider fetches the registered provider from the configured name
func fetchProvider() (Provider, error) {

	p, ok := providers[config.Viper().GetString("provider")]
	if !ok {
		return nil, errors.New("invalid provider")
	}

	return p, nil
}
