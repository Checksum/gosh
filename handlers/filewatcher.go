package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Checksum/notify"
)

// FileWatcher recursively watches the served directory for changes
type FileWatcher struct {
	Handler    http.Handler
	Broker     *Broker
	Watcher    chan notify.EventInfo
	Extensions []string
}

// NewFileWatcher returns a new file watcher
func NewFileWatcher(dir, ext string, h http.Handler) *FileWatcher {
	b := NewBroker()
	c := make(chan notify.EventInfo)
	// Start the broker
	b.Start()
	// Start the listener
	err := notify.WatchWithFilter(dir+"/...", c, filterPath, notify.Create|notify.Write)

	if err != nil {
		log.Fatal(err)
	}

	fw := &FileWatcher{h, b, c, strings.Split(ext, ",")}

	go func() {
		for {
			select {
			case event := <-c:
				ext := filepath.Ext(event.Path())
				if fw.isValidExt(ext) {
					log.Println("File changed, reloading")
					b.Messages <- "reload:" + event.Path()
				}
			}
		}
	}()

	return fw

}

func (fw *FileWatcher) isValidExt(ext string) bool {
	if ext != "" && len(ext) > 1 {
		ext = ext[1:]
		for _, val := range fw.Extensions {
			if val == ext {
				return true
			}
		}
	}
	return false
}

// This doesn't seem to be working!!
func filterPath(path string) bool {
	return true
}
