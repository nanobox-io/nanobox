//
package notify

import (
	"github.com/go-fsnotify/fsnotify"
)

type (
	notify struct{}
	Notify interface {
		Watch(path string, handle func(e *fsnotify.Event) error) error
	}
)

var (
	Default Notify = notify{}
)

func (notify) Watch(path string, handle func(e *fsnotify.Event) error) error {
	return Watch(path, handle)
}
