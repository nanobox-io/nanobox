// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package server

type (
	server struct{}
	Server interface {
		Bootstrap(params string) error
		Build(params string) error
		Deploy(params string) error
		Exec(kind, params string) error
		IsContainerExec(args []string) (found bool)
		NotifyRebuild(event *fsnotify.Event) error
		NotifyServer(event *fsnotify.Event) error
		Lock() error
		Unlock()
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

func (vagrant) Bootstrap(params string) error {
	return Bootstrap(params)
}

func (vagrant) Build(params string) error {
	return Build(params)
}

func (vagrant) Deploy(params string) error {
	return Deploy(params)
}

func (vagrant) Exec(kind, params string) error {
	return Exec(kind, params)
}

func (vagrant) IsContainerExec(args []string) (found bool) {
	return IsContainerExec(args)
}

func (vagrant) NotifyRebuild(event *fsnotify.Event) error {
	return NotifyRebuild(event)
}

func (vagrant) NotifyServer(event *fsnotify.Event) error {
	return NotifyServer(event)
}

func (vagrant) Lock() error {
	return Lock()
}

func (vagrant) Unlock() {
	Unlock()
}

func (vagrant) Logs(params string) error {
	return Logs(params)
}

func (vagrant) Ping() (bool, error) {
	return Ping()
}

func (vagrant) Get(path string, v interface{}) (*http.Response, error) {
	return Get(path, v)
}

func (vagrant) Post(path, contentType string, body io.Reader) (*http.Response, error) {
	return Post(path, contentType, body)
}

func (vagrant) Put(path string, body io.Reader) (*http.Response, error) {
	return Put(path, body)
}

func (vagrant) Suspend() error {
	return Suspend()
}

func (vagrant) Update(params string) error {
	return Update(params)
}
