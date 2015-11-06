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
		Tar(path string, writers ...io.Writer) error
		Untar(path string, r io.Reader)
		Download(path string, w io.Writer) error
		Progress(path string, w io.Writer) error
	}
)

var (
	Default File = file{}
)

func (file) Tar(path string, writers ...io.Writer) error {
	return Tar(path, writers...)
}

func (file) Untar(path string, r io.Reader) {
	Untar(path, r)
}

func (file) Download(path string, w io.Writer) error {
	return Download(path, w)
}

func (file) Progress(path string, w io.Writer) error {
	return Progress(path, w)
}
