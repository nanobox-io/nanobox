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
		id         string
		downloadCurrent int
		downloadTotal int
		downloaded int
		extracted  int
	}

	// DockerPercentDisplay ...
	DockerPercentDisplay struct {
		Output   io.Writer
		Prefix   string
		parts    []*DockerPercentPart
		leftover []byte
	}
)

// update ...
func (part *DockerPercentPart) update(status Status) {
	switch status.Status {

	//
	case "Downloading":
		part.downloadCurrent = status.Details.Current
		part.downloadTotal = status.Details.Total
		part.downloaded = int(float64(part.downloadCurrent) / float64(part.downloadTotal) * 100.0)

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

	//
	for _, v := range display.parts {
		count++

		if v.downloaded != 100 {
			return fmt.Sprintf("Downloading layer %2d/%d: %2d%% %s", count, len(display.parts), v.downloaded, displaySize(v))
		} else if count == len(display.parts) {
			return fmt.Sprintf("Extracting layer %2d/%d: %2d%%", count, len(display.parts), v.extracted)
		}
	}

	return ""
}

// Write ...
func (display *DockerPercentDisplay) Write(data []byte) (int, error) {
	// set it if not set already
	if display.parts == nil {
		display.parts = []*DockerPercentPart{}
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
			found := false

			for _, part := range display.parts {
				if part.id == status.ID {
					part.update(status)
					found = true
					break
				}
			}

			if !found {
				part := &DockerPercentPart{id: status.ID}
				part.update(status)
				display.parts = append(display.parts, part)
			}
		}

		fmt.Fprintf(display.Output, "\r\x1b[K")
		fmt.Fprintf(display.Output, "%s %s", display.Prefix, display.show())

		if strings.HasPrefix(status.Status, "Status:") {
			// maybe we want to display the status line here
		}
	}

	return len(data), nil
}

func displaySize(part *DockerPercentPart) string {
	switch {
	case part.downloadTotal > 1024*1024:
		// mbps
		return fmt.Sprintf("%.2fMB/%.2fMB", float64(part.downloadCurrent)  / float64(1024*1024), float64(part.downloadTotal)  / float64(1024*1024))
	case part.downloadTotal > 1024:
		// kbps
		return fmt.Sprintf("%.2fKB/%.2fKB", float64(part.downloadCurrent)  / float64(1024), float64(part.downloadTotal)  / float64(1024))
	default:
		// bps
		return fmt.Sprintf("%.2fB/%.2fB", float64(part.downloadCurrent), float64(part.downloadTotal))
	}
}
