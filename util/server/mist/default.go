// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package mist

type (
	mist struct{}
	Mist interface {
		Listen(tags []string, handle func(string) error) error
		Stream(tags []string, handle func(Log))
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
