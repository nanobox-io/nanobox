//
package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// "Pulling fs layer"
// "Waiting"
// "Downloading"
// "Verifying Checksum"
// "Download complete"
// "Extracting"
// "Pull complete"
// "Status: **"

const (
	Pulling = iota
	Waiting
	Downloading
	Verifying
	DLComplete
	Extracting
	Complete
)

var (
	displays = map[int]string{
		Pulling:     "_",
		Waiting:     ".",
		Downloading: "_",
		Verifying:   "^",
		DLComplete:  "↓",
		Extracting:  "/",
		Complete:    "✓",
	}
	nextDownloader = map[string]string{
		"":  "_", // catch any messups during a download
		"_": "-",
		"-": "‾",
		"‾": "_",
	}
	nextExtracter = map[string]string{
		"":   "/", // catch messups during extracting
		"/":  "\\",
		"\\": "/",
	}
)

type Status struct {
	Status string `json:"status,omitempty"`
	ID     string `json:"id,omitempty"`
}

type DockerImageDisplaySimple struct {
	Prefix     string
	idDisplays map[string]string
	leftover   []byte
}

func (self *DockerImageDisplaySimple) bar() string {
	// order them
	keys := []string{}
	for k, _ := range self.idDisplays {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))

	s := ""
	for _, key := range keys {
		s = s + self.idDisplays[key]
	}
	return fmt.Sprintf("[%s]", s)
}

func (self *DockerImageDisplaySimple) Write(data []byte) (int, error) {
	// set it if not set already
	if self.idDisplays == nil {
		self.idDisplays = map[string]string{}
	}
	// create a buffer with the old leftovers and the new data
	buffer := bytes.NewBuffer(append(self.leftover, data...))
	// clear out the leftovers
	self.leftover = []byte{}

	for {
		line, err := buffer.ReadBytes('\n')
		if err == io.EOF {
			self.leftover = line
			break
		}
		// take the line and turn it into a status
		status := Status{}
		json.Unmarshal(line, &status)
		if err != nil {
			fmt.Println("what?", err)
		}

		switch status.Status {
		case "Pulling fs layer":
			self.idDisplays[status.ID] = displays[Pulling]
		case "Waiting":
			self.idDisplays[status.ID] = displays[Waiting]
		case "Downloading":
			prev := self.idDisplays[status.ID]
			if prev == "_" || prev == "-" || prev == "‾" {
				self.idDisplays[status.ID] = nextDownloader[prev]
			} else {
				self.idDisplays[status.ID] = displays[Downloading]
			}
		case "Verifying Checksum":
			self.idDisplays[status.ID] = displays[Verifying]
		case "Download complete":
			self.idDisplays[status.ID] = displays[DLComplete]
		case "Extracting":
			prev := self.idDisplays[status.ID]
			if prev == "/" || prev == "\\" {
				self.idDisplays[status.ID] = nextExtracter[prev]
			} else {
				self.idDisplays[status.ID] = displays[Extracting]
			}
		case "Pull complete":
			self.idDisplays[status.ID] = displays[Complete]
		default:
			if bytes.HasPrefix([]byte(status.Status), []byte("Status:")) {
				// clear the downloader
				fmt.Fprintf(os.Stdout, "%c[2K\r", 27)
				// fmt.Printf("\n%s\n", status.Status)
				return len(data), nil
			}
			// fmt.Printf("couldnt find athing for : %+v", status)
		}
		fmt.Fprintf(os.Stdout, "%c[2K\r", 27)
		fmt.Printf("%s: %s", self.Prefix, self.bar())

	}

	return len(data), nil
}
