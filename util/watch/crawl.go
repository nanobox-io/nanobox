package watch

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type crawl struct {
	path string

	events chan event
	done   chan struct{}

	started bool
	files   map[string]time.Time
}

func newCrawlWatcher(path string) Watcher {
	return &crawl{
		path:   path,
		events: make(chan event, 10),
		files:  map[string]time.Time{},
	}
}

func (c *crawl) watch() error {
	// fill in the files list
	if err := c.populateFiles(); err != nil {
		return err
	}
	// start the continual watcher
	go c.start()

	return nil
}

// retrieve the event channel
func (c *crawl) eventChan() chan event {
	return c.events
}

func (c *crawl) close() error {
	close(c.done)
	return nil
}

func (c *crawl) populateFiles() error {
	return filepath.Walk(c.path, c.walkFunc)
}

// add a file that is being walked to the watch system
func (c *crawl) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	for _, ignoreName := range ignoreFile {
		if strings.HasSuffix(path, ignoreName) {
			if info.IsDir() {
				// if the thing we are ignoring is a directory
				return filepath.SkipDir
			}
			// if its not just skip the file
			return nil
		}
	}

	// read the file with the md5 library and generate a hash
	if !info.IsDir() {
		val, ok := c.files[path]
		if c.started && (!ok || !val.Equal(info.ModTime().Round(5*time.Second))) {
			// this is a new file or the file has been changed
			c.events <- event{file: path}
		}

		// update my cached files
		// the rounding is so we dont detect the change that we make
		c.files[path] = info.ModTime().Round(5 * time.Second)
	}

	return nil
}

func (c *crawl) start() {
	c.started = true
	for {
		select {
		// sleep for a second between walking the tree
		// this could be made variable
		case <-time.After(time.Second):
			err := filepath.Walk(c.path, c.walkFunc)
			if err != nil {
				c.events <- event{error: err}
				close(c.events)
				return
			}

			// if we are asked to close then close grace fully
		case <-c.done:
			close(c.events)
			return
		}
	}
}
