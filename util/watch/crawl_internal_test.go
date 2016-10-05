package watch

import (
	"os"
	"os/exec"
	"testing"
)

func TestCrawlFiles(t *testing.T) {
	os.MkdirAll("/tmp/nanobox/", 0777)
	crawlWatcher := newCrawlWatcher("/tmp/nanobox/")
	err := crawlWatcher.watch()
	if err != nil {
		t.Fatalf("failed to watch: %s", err)
	}
	defer crawlWatcher.close()

	exec.Command("touch", "/tmp/nanobox/crawl.tmp").Run()

	// pull the first event off the channel
	ev := <-crawlWatcher.eventChan()

	if ev.file != "/tmp/nanobox/crawl.tmp" {
		t.Errorf("the wrong file path came out %s", ev.file)
	}
	if ev.error != nil {
		t.Errorf("an error occurred %s", ev.error)
	}
}
