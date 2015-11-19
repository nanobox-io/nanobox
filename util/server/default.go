//
package server

import (
	"io"
	"net/http"

	"github.com/go-fsnotify/fsnotify"
)

type (
	server struct{}
	Server interface {
		Bootstrap(params string) error
		Build(params string) error
		Deploy(params string) error
		Exec(where, params string) error
		IsContainerExec(args []string) (found bool)
		NotifyRebuild(event *fsnotify.Event) error
		NotifyServer(event *fsnotify.Event) error
		Lock()
		Unlock()
		NewLogger(path string)
		Logs(params string) error
		Ping() (bool, error)
		Get(path string, v interface{}) (*http.Response, error)
		Post(path, contentType string, body io.Reader) (*http.Response, error)
		Put(path string, body io.Reader) (*http.Response, error)
		Suspend() error
		Update(params string) error
	}
)

var (
	Default Server = server{}
)

func (server) Bootstrap(params string) error {
	return Bootstrap(params)
}

func (server) Build(params string) error {
	return Build(params)
}

func (server) Deploy(params string) error {
	return Deploy(params)
}

func (server) Exec(where, params string) error {
	return Exec(where, params)
}

func (server) IsContainerExec(args []string) (found bool) {
	return IsContainerExec(args)
}

func (server) NotifyRebuild(event *fsnotify.Event) error {
	return NotifyRebuild(event)
}

func (server) NotifyServer(event *fsnotify.Event) error {
	return NotifyServer(event)
}

func (server) Lock() {
	Lock()
}

func (server) Unlock() {
	Unlock()
}

func (server) NewLogger(path string) {
	NewLogger(path)
}

func (server) Logs(params string) error {
	return Logs(params)
}

func (server) Ping() (bool, error) {
	return Ping()
}

func (server) Get(path string, v interface{}) (*http.Response, error) {
	return Get(path, v)
}

func (server) Post(path, contentType string, body io.Reader) (*http.Response, error) {
	return Post(path, contentType, body)
}

func (server) Put(path string, body io.Reader) (*http.Response, error) {
	return Put(path, body)
}

func (server) Suspend() error {
	return Suspend()
}

func (server) Update(params string) error {
	return Update(params)
}
