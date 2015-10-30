// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package file

import "io"

type (
	file struct{}
	File interface {
		Gzip()
		Tar()
		TarBall()
		Download(path string, w io.Writer) error
		Progress(path string, w io.Writer) error
	}
)

var (
	Default File = file{}
)

func (file) Gzip() {
	Gzip()
}

func (file) Tar() {
	Tar()
}

func (file) TarBall() {
	TarBall()
}

func (file) Download(path string, w io.Writer) error {
	return Download(path, w)
}

func (file) Progress(path string, w io.Writer) error {
	return Progress(path, w)
}
