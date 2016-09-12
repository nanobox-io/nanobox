package fileutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Download a file from a url to the path specified
func Download(url, path string) error {
	// create the file at the path specified
	fd, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file to download into: %s", err.Error())
	}

	// ensure the file descriptor is closed
	defer fd.Close()

	// fetch the file from the url specified
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch file to download: %s", err.Error())
	}

	// ensure the body is closed
	defer res.Body.Close()

	// create a buffer to read into
	b := make([]byte, 2048)

	for {
		// read the response body (streaming)
		n, err := res.Body.Read(b)

		// write the contents of our buffer to the file
		fd.Write(b[:n])

		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("failed to read body: %s", err.Error())
		}
	}

	return nil
}
