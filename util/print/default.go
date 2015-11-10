//
package print

type (
	print struct{}
	Print interface {
		Verbose(msg string)
		Silence(msg string)
		Color(msg string, v ...interface{})
		Prompt(p string, v ...interface{}) string
		Password(p string) string
	}
)

var (
	Default Print = print{}
)

func (print) Verbose(msg string) {
	Verbose(msg)
}

func (print) Silence(msg string) {
	Silence(msg)
}

func (print) Color(msg string, v ...interface{}) {
	Color(msg, v...)
}

func (print) Prompt(p string, v ...interface{}) string {
	return Prompt(p, v)
}

func (print) Password(p string) string {
	return Password(p)
}
