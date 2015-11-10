//
package vagrant

type (
	vagrant struct{}
	Vagrant interface {
		HaveImage() bool
		Install() error
		Update() error
		Destroy() error
		Init()
		NewLogger(path string)
		Reload() error
		Resume() error
		SSH() error
		Status() string
		Suspend() error
		Up() error
	}
)

var (
	Default Vagrant = vagrant{}
)

func (vagrant) Up() error {
	return Up()
}

func (vagrant) Suspend() error {
	return Suspend()
}

func (vagrant) Status() (status string) {
	return Status()
}

func (vagrant) SSH() error {
	return SSH()
}

func (vagrant) NewLogger(path string) {
	NewLogger(path)
}

func (vagrant) Reload() error {
	return Reload()
}

func (vagrant) Resume() error {
	return Resume()
}

func (vagrant) Init() {
	Init()
}

func (vagrant) Destroy() error {
	return Destroy()
}

func (vagrant) HaveImage() bool {
	return HaveImage()
}

func (vagrant) Install() error {
	return Install()
}

func (vagrant) Update() error {
	return Update()
}
