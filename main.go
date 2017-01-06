package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Checksum/gosh/handlers"
	"github.com/toqueteos/webbrowser"

	"github.com/Zillode/notify"
	gorilla "github.com/gorilla/handlers"
	"github.com/pkg/profile"
	"github.com/zenazn/goji/web/middleware"
)

func main() {
	var port, dir, ext string
	var noCache, watch, openBrowser, spa, cpuprofile bool

	flag.StringVar(&port, "port", "8000", "The port to bind to")
	flag.StringVar(&dir, "dir", ".", "The directory to serve")
	flag.BoolVar(&noCache, "no-cache", true, "Disable HTTP caching")
	flag.BoolVar(&watch, "watch", false, "Reload browser on file change")
	flag.BoolVar(&openBrowser, "open", true, "Open a browser to serve the page")
	flag.BoolVar(&spa, "spa", true, "Serve a single page application. All unmatched routes are forwarded to index.html")
	flag.StringVar(&ext, "ext", "js,css,html", "Comma separated file extensions to watch for change")
	flag.BoolVar(&cpuprofile, "cpuprofile", false, "Record profile for debugging")
	flag.Parse()

	if cpuprofile {
		defer profile.Start().Stop()
	}

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

	if spa {
		handler = handlers.NewFileServer(handler)
	}

	// Handle all other paths
	mux.Handle("/", handler)

	if openBrowser {
		webbrowser.Open("http://localhost:" + port)
	}
	log.Printf("Serving %s on port %s", dir, port)
	log.Fatal(http.ListenAndServe("localhost:"+port, mux))
}
