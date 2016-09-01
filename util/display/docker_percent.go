package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type (

	// Status ...
	// {"status":"Downloading","progressDetail":{"current":676,"total":755},"progress":"[============================================\u003e      ]    676 B/755 B","id":"166102ec41af"}
	Status struct {
		Status  string  `json:"status,omitempty"`
		ID      string  `json:"id,omitempty"`
		Details Details `json:"progressDetail"`
	}

	// Details ...
	Details struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	}

	// DockerPercentPart ...
	DockerPercentPart struct {
		downloaded int
		extracted  int
	}

	// DockerPercentDisplay ...
	DockerPercentDisplay struct {
		Output   io.Writer
		Prefix   string
		parts    map[string]DockerPercentPart
		leftover []byte
	}
)

// update ...
func (part *DockerPercentPart) update(status Status) {
	switch status.Status {

	//
	case "Downloading":
		current := status.Details.Current
		total := status.Details.Total
		part.downloaded = int(float64(current) / float64(total) * 100.0)

	//
	case "Download complete":
		part.downloaded = 100

	//
	case "Extracting":
		current := status.Details.Current
		total := status.Details.Total
		part.extracted = int(float64(current) / float64(total) * 100.0)

	//
	case "Pull complete":
		part.extracted = 100

	//
	case "Already exists":
		part.downloaded = 100
		part.extracted = 100

	//
	default:
		// there is a chance if given a tag (nanobox/build:v1)
		// it will be able to pull a part from the non labeled parts
		if strings.HasPrefix(status.Status, "Pulling from") {
			part.downloaded = 100
			part.extracted = 100
		}
	}
}

// show ...
func (display *DockerPercentDisplay) show() string {

	// order them
	count := 0
	downloaded := 0
	extracted := 0

	//
	for _, v := range display.parts {
		count++
		downloaded += v.downloaded
		extracted += v.extracted
	}

	//
	if count == 0 {
		count = 1
	}

	//
	return fmt.Sprintf("Downloaded: %3d%% Extracted: %3d%% Total: %3d%%", downloaded/count, extracted/count, (downloaded/count+extracted/count)/2)
}

// Write ...
func (display *DockerPercentDisplay) Write(data []byte) (int, error) {
	// set it if not set already
	if display.parts == nil {
		display.parts = map[string]DockerPercentPart{}
	}
	// create a buffer with the old leftovers and the new data
	buffer := bytes.NewBuffer(append(display.leftover, data...))
	// clear out the leftovers
	display.leftover = []byte{}

	for {
		line, err := buffer.ReadBytes('\n')
		if err == io.EOF {
			display.leftover = line
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
			part, ok := display.parts[status.ID]
			if !ok {
				part = DockerPercentPart{}
			}
			part.update(status)
			display.parts[status.ID] = part
		}
		// fmt.Fprintf(display.Output, "%c[2K\r", 27)
		fmt.Fprintf(display.Output, "%s %s", display.Prefix, display.show())
		if strings.HasPrefix(status.Status, "Status:") {
			// maybe we want to display the status line here
		}
	}

	return len(data), nil
}
