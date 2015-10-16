// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package vagrant

type (
	vagrant struct{}
	Vagrant interface {
		HaveImage() bool
		Install() error
		Update() error
		Destroy() error
		Init() error
		Reload() error
		SSH() error
		Status() string
		Suspend() error
		Up() error
		Run() (err error)
	}
)

var (
	Default Vagrant = vagrant{}
)

func (vagrant) Run() (err error) {
	return Run()
}

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

func (vagrant) Reload() error {
	return Reload()
}

func (vagrant) Init() error {
	return Init()
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
