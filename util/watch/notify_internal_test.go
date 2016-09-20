package watch

import (
	"testing"
	"io/ioutil"
	"os"
	"time"
)

func TestNotifyFiles(t *testing.T) {
	os.MkdirAll("/tmp/nanobox/", 0777)
	notifyWatcher := newNotifyWatcher("/tmp/nanobox/")
	err := notifyWatcher.watch()
	if err != nil {
		t.Fatalf("failed to watch: %s", err)
	}

	defer notifyWatcher.close()
	<-time.After(time.Second)
	ioutil.WriteFile("/tmp/nanobox/notify.tmp", []byte("hi"), 0777)

	// pull the first event off the channel
	ev := <-notifyWatcher.eventChan()

	if ev.file != "/tmp/nanobox/notify.tmp" {
		t.Errorf("the wrong file path came out %s", ev.file)
	}
	if ev.error != nil {
		t.Errorf("an error occurred %s", ev.error)
	}
}
