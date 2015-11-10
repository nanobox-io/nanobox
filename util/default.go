//
package util

type (
	util struct{}
	Util interface {
		MD5sMatch(string, string) (bool, error)
	}
)

var (
	Default Util = util{}
)

func (_ util) MD5sMatch(localPath, remotePath string) (bool, error) {
	return MD5sMatch(localPath, remotePath)
}
