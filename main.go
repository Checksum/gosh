package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Checksum/gosh/handlers"
	"github.com/Zillode/notify"

	gorilla "github.com/gorilla/handlers"
	"github.com/zenazn/goji/web/middleware"
)

func main() {
	var port, dir, ext string
	var noCache, watch bool

	flag.StringVar(&port, "port", "8000", "The port to bind to")
	flag.StringVar(&dir, "dir", ".", "The directory to serve")
	flag.BoolVar(&noCache, "no-cache", false, "Disable HTTP caching")
	flag.BoolVar(&watch, "watch", false, "Reload browser on file change")
	flag.StringVar(&ext, "ext", "js,css,html", "Comma separated file extensions to watch for change")
	flag.Parse()

	mux := http.NewServeMux()
	handler := gorilla.LoggingHandler(os.Stdout, http.FileServer(http.Dir(dir)))

	if noCache {
		handler = middleware.NoCache(handler)
	}

	// If we have to watch the file for changes, pass the request through
	// the appropriate handlers.
	if watch {
		// This is to intercept the original response to inject our javascript
		handler = handlers.NewInjectingHandler(handler)
		// Create a new watcher that will notify the page through Server-Sent Events
		watchHandler := handlers.NewFileWatcher(dir, ext, handler)
		defer notify.Stop(watchHandler.Watcher)
		// This endpoint is serviced by the SSE broker, which receives a signal
		// when a file system change occurs, and sends out a command to the browser
		// to reload the page (see also: script.go)
		mux.Handle("/__events", watchHandler.Broker)
		log.Println("Watching for changes..")
	}

	// Handle all other paths
	mux.Handle("/", handler)

	log.Printf("Serving %s on port %s", dir, port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
