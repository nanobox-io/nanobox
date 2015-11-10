//
package s3

import (
	"io"
	"net/http"
)

type (
	s3 struct{}
	S3 interface {
		RequestURL(path string) (string, error)
		Download(path string) (*http.Response, error)
		Upload(path string, body io.Reader) error
		Request(method, path string, body io.Reader) (*http.Response, error)
	}
)

var (
	Default S3 = s3{}
)

func (s3) RequestURL(path string) (string, error) {
	return RequestURL(path)
}

func (s3) Download(path string) (*http.Response, error) {
	return Download(path)
}

func (s3) Upload(path string, body io.Reader) error {
	return Upload(path, body)
}

func (s3) Request(method, path string, body io.Reader) (*http.Response, error) {
	return Request(method, path, body)
}
