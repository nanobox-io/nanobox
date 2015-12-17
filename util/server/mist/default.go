//
package mist

import mistClient "github.com/nanopack/mist/core"

type (
	mist struct{}
	Mist interface {
		Listen(tags []string, handle func(string) error) error
		Stream(tags []string, handle func(Log))
		Connect(tags []string, handle func(Log)) (client mistClient.Client, err error)
		ProcessLog(log Log)
		DeployUpdates(status string) (err error)
		BuildUpdates(status string) (err error)
		BootstrapUpdates(status string) (err error)
		ImageUpdates(status string) (err error)
		PrintLogStream(log Log)
		ProcessLogStream(log Log)
	}
)

var (
	Default Mist = mist{}
)

func (mist) Listen(tags []string, handle func(string) error) error {
	return Listen(tags, handle)
}

func (mist) Stream(tags []string, handle func(Log)) {
	Stream(tags, handle)
}

func (mist) Connect(tags []string, handle func(Log)) (clinet mistClient.Client, err error) {
	return Connect(tags, handle)
}

func (mist) ProcessLog(log Log) {
	ProcessLog(log)
}

func (mist) DeployUpdates(status string) (err error) {
	return DeployUpdates(status)
}

func (mist) BuildUpdates(status string) (err error) {
	return BuildUpdates(status)
}

func (mist) BootstrapUpdates(status string) (err error) {
	return BootstrapUpdates(status)
}

func (mist) ImageUpdates(status string) (err error) {
	return ImageUpdates(status)
}

func (mist) PrintLogStream(log Log) {
	PrintLogStream(log)
}

func (mist) ProcessLogStream(log Log) {
	ProcessLogStream(log)
}
