package display

import (
	"io"	
	"strings"
	"os"
	"fmt"
)

const bytesPerMB = 1024 * 1024

type DownloadPercent struct {
	current int64
	Total   int64
	Output  io.Writer
}

func (dp *DownloadPercent) Copy(writer io.Writer, reader io.Reader) (err error) {
	// initialize variables
	buf := make([]byte, 32*1024)
	dp.current = 0

	if dp.Output == nil {
		dp.Output = os.Stdout
	}

	for {
		dp.UpdateDisplay()
		readBytes, readErr := reader.Read(buf)
		if readBytes > 0 {
			writtenBytes, writeErr := writer.Write(buf[0:readBytes])
			if writeErr != nil {
				err = writeErr
				break
			}
			if readBytes != writtenBytes {
				err = io.ErrShortWrite
				break
			}

		}

		// success
		if readErr == io.EOF {
			break
		}

		// failure
		if readErr != nil {
			err = readErr
			break
		}

		dp.current += int64(readBytes)
		dp.UpdateDisplay()
	}

	if err == nil {
		dp.current = dp.Total
		dp.UpdateDisplay()
	}
	return err
}

func (dp *DownloadPercent) UpdateDisplay() {
	// clear the link
	fmt.Fprintf(dp.Output, "\r\x1b[K")

	if dp.Total == 0 {
		dp.SimpleDisplay()
		return
	}

	// show download progress: 0.0/0.0MB [*** progress *** 0.0%]
	currentInMB := float64(dp.current)/bytesPerMB
	totalInMB := float64(dp.Total)/bytesPerMB
	percent := (float64(dp.current) / float64(dp.Total)) * 100

	fmt.Fprintf(dp.Output, "\r   %.2f/%.2fMB [%-41s %.2f%%]", currentInMB, totalInMB, strings.Repeat("*", int(percent/2.5)), percent)

}

func (dp *DownloadPercent) SimpleDisplay() {
	fmt.Fprintf(dp.Output, "   %.2fMB", float64(dp.current)/bytesPerMB)
}
