package server

import (
	"testing"
	"time"
	"io"
)

	
func TestTimeoutReader(t *testing.T) {
	// create a new timeout reader
	timoutReader = &TimeoutReader{
		Files:   make(chan string, 1),
		timeout: time.Second,
	}

	bytes := make([]byte, 4)

	timoutReader.Files <- "1"
	n, err := timoutReader.Read(bytes)
	if n != 2 || err != nil || string(bytes[:n]) != "1\n" {
		t.Errorf("the reader didnt output the right data n %d, err %+v, data %q", n, err, bytes)
	}
	
	timoutReader.Files <- "123"
	n, err = timoutReader.Read(bytes)
	if n != 4 || err != nil || string(bytes[:n]) != "123\n" {
		t.Errorf("the reader didnt output the right data n %d, err %+v, data %q", n, err, bytes)
	}

	// now attempt a read without writing
	// should timeout
	n, err = timoutReader.Read(bytes)
	if n != 0 || err != io.EOF {
		t.Errorf("the timeout reader didnt timeout")
	}
}