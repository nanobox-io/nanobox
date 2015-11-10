//
package auth

type (
	auth struct{}
	Auth interface {
		Authenticate() (string, string)
		Reauthenticate() (string, string)
	}
)

var (
	Default Auth = auth{}
)

//
func (auth) Authenticate() (string, string) {
	return Authenticate()
}

//
func (auth) Reauthenticate() (string, string) {
	return Reauthenticate()
}
