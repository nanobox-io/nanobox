// Package util ...
package util

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/mitchellh/colorstring"
)

// MD5sMatch determines if a local MD5 matches a remote MD5
func MD5sMatch(localFile, remotePath string) (bool, error) {

	// read the local file; will return os.PathError if doesn't exist
	b, err := ioutil.ReadFile(localFile)
	if err != nil {
		return false, err
	}

	// get local md5 checksum (as a string)
	localMD5 := fmt.Sprintf("%x", md5.Sum(b))

	// GET remote md5
	res, err := http.Get(remotePath)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	// read the remote md5 checksum
	remoteMD5, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// compare checksum's
	return strings.TrimSpace(localMD5) == strings.TrimSpace(string(remoteMD5)), nil
}

// Progress downloads a file with a fancy progress bar
func Progress(path string, w io.Writer) error {

	//
	download, err := http.Get(path)
	if err != nil {
		return err
	}
	defer download.Body.Close()

	var percent float64
	var down int

	// format the response content length to be more 'friendly'
	total := float64(download.ContentLength) / math.Pow(1024, 2)

	// create a 'buffer' to read into
	p := make([]byte, 2048)

	//
	for {

		// read the response body (streaming)
		n, err := download.Body.Read(p)

		// write to our buffer
		w.Write(p[:n])

		// update the total bytes read
		down += n

		switch {
		default:
			// update the percent downloaded
			percent = (float64(down) / float64(download.ContentLength)) * 100

			// show download progress: 0.0/0.0MB [*** progress *** 0.0%]
			fmt.Printf("\r   %.2f/%.2fMB [%-41s %.2f%%]", float64(down)/math.Pow(1024, 2), total, strings.Repeat("*", int(percent/2.5)), percent)

		// if no more files are found return
		case download.ContentLength < 1:
			fmt.Printf("\r   %.2fMB", float64(down)/math.Pow(1024, 2))
		}

		// detect EOF and break the 'stream'
		if err != nil {
			if err == io.EOF {
				fmt.Println("")
				break
			} else {
				return err
			}
		}
	}

	return nil
}

// Prompt will prompt for input from the shell and return a trimmed response
func Prompt(p string, v ...interface{}) string {
	reader := bufio.NewReader(os.Stdin)

	//
	fmt.Print(colorstring.Color(fmt.Sprintf(p, v...)))

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		// config.Fatal("[util/print] reader.ReadString() failed", err.Error())
	}

	return strings.TrimSpace(input)
}
