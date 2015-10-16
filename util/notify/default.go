// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

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
