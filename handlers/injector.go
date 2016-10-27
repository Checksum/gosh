package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// https://justinas.org/writing-http-middleware-in-go/

// TransformationFunc is the function used to transform the original response
type TransformationFunc func(r io.Reader) (string, error)

// InjectMiddleware intercepts the response and
// modifies it in some way
type InjectMiddleware struct {
	handler   http.Handler
	transform TransformationFunc
}

func (m *InjectMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Intercepts only requests with text/html in the Accept request header
	// We need to do this so that we inject the script only into html pages
	// TODO: Prevent multiple injections for same page (ex: embedded pages, xhr templates, etc)
	if !strings.Contains(r.Header.Get("Accept"), "text/html") {
		m.handler.ServeHTTP(w, r)
		return
	}

	rec := httptest.NewRecorder()
	// passing a ResponseRecorder instead of the original RW
	m.handler.ServeHTTP(rec, r)
	// Original response body
	body := rec.Body.Bytes()

	// Modify response only if a valid HTTP response
	if rec.Code >= 200 || rec.Code < 300 {
		// If we get here, we have to intercept the request and inject our script
		log.Println("Intercepting: " + r.URL.String())

		content, err := m.transform(bytes.NewReader(body))

		if err == nil && content != "" {
			body = []byte(content)
		}
	}

	// we copy the original headers first
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	// Set the correct Content-Length now that we are injecting
	// extra data to the response
	contentLength := len(body)
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.WriteHeader(rec.Code)

	// Write out the modified response
	w.Write(body)
}

// NewInjectingHandler returns a http.Handler that intercepts
// the original response and modifies it
func NewInjectingHandler(h http.Handler) http.Handler {
	return &InjectMiddleware{h, transformResponse}
}

func transformResponse(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return "", err
	}
	// Even if body is empty, goquery inserts the missing tags
	// automatically when calling doc.Html()
	body := doc.Find("body").First()
	body.AppendHtml(ScriptContent)

	return doc.Html()
}
