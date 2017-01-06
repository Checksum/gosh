package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

// SPAMiddleware intercepts the response and
// modifies it in some way
type SPAMiddleware struct {
	handler http.Handler
}

func (m *SPAMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(".", r.URL.Path)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		http.ServeFile(w, r, path)
	} else {
		http.ServeFile(w, r, filepath.Join(".", "index.html"))
	}
}

// NewFileServer creates a new middleware for single page applications
func NewFileServer(h http.Handler) http.Handler {
	return &SPAMiddleware{h}
}
