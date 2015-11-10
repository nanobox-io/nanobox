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
