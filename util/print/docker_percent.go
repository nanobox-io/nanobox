//
package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
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

type DockerPercentPart struct {
	downloaded int
	extracted  int
}

func (self *DockerPercentPart) update(status Status) {
	switch status.Status {
	case "Downloading":
		current := status.Details.Current
		total := status.Details.Total
		self.downloaded = int(float64(current) / float64(total) * 100.0)
	case "Download complete":
		self.downloaded = 100
	case "Extracting":
		current := status.Details.Current
		total := status.Details.Total
		self.extracted = int(float64(current) / float64(total) * 100.0)
	case "Pull complete":
		self.extracted = 100
	case "Already exists":
		self.downloaded = 100
		self.extracted = 100
	default:
		// there is a chance if given a tag (nanobox/build:v1)
		// it will be able to pull a part from the non labeled parts
		if strings.HasPrefix(status.Status, "Pulling from") {
			self.downloaded = 100
			self.extracted = 100
		}
	}
}

// {"status":"Downloading","progressDetail":{"current":676,"total":755},"progress":"[============================================\u003e      ]    676 B/755 B","id":"166102ec41af"}
type Status struct {
	Status  string  `json:"status,omitempty"`
	ID      string  `json:"id,omitempty"`
	Details Details `json:"progressDetail"`
}

type Details struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type DockerPercentDisplay struct {
	Prefix   string
	parts    map[string]DockerPercentPart
	leftover []byte
}

func (self *DockerPercentDisplay) show() string {
	// order them
	count := 0
	downloaded := 0
	extracted := 0
	for _, v := range self.parts {
		count++
		downloaded += v.downloaded
		extracted += v.extracted
	}
	if count == 0 {
		count = 1
	}
	return fmt.Sprintf("Downloaded: %3d%% Extracted: %3d%% Total: %3d%%", downloaded/count, extracted/count, (downloaded/count+extracted/count)/2)
}

func (self *DockerPercentDisplay) Write(data []byte) (int, error) {
	// set it if not set already
	if self.parts == nil {
		self.parts = map[string]DockerPercentPart{}
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
			fmt.Println(err)
			continue
		}
		if status.ID != "latest" && status.ID != "" {
			part, ok := self.parts[status.ID]
			if !ok {
				part = DockerPercentPart{}
			}
			part.update(status)
			self.parts[status.ID] = part
		}
		fmt.Fprintf(os.Stdout, "%c[2K\r", 27)
		fmt.Printf("%s %s", self.Prefix, self.show())
		if strings.HasPrefix(status.Status, "Status:") {
			fmt.Printf("\n")
			// maybe we want to display the status line here
		}
	}

	return len(data), nil
}
