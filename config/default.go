//
package config

var (
	Default Config = config{}
)

type (
	Config interface {
		Fatal(string, string)
		Root() string
		ParseConfig(path string, v interface{}) error
		Debug(msg string, debug bool)
		Info(msg string)
		Error(msg, err string)
		ParseNanofile() *NanofileConfig
		ParseVMfile() *VMfileConfig
	}

	config struct {
	}
)

func (config) ParseVMfile() *VMfileConfig {
	return ParseVMfile()
}

func (config) ParseNanofile() *NanofileConfig {
	return ParseNanofile()
}

func (config) Debug(msg string, debug bool) {
	Debug(msg, debug)
}

func (config) Info(msg string) {
	Info(msg)
}

func (config) Error(msg, err string) {
	Error(msg, err)
}

func (config) ParseConfig(path string, v interface{}) error {
	return ParseConfig(path, v)
}

func (config) Fatal(msg, err string) {
	Fatal(msg, err)
}

func (config) Root() string {
	return Root
}
