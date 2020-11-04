package app

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"one-off-email/models"
	"time"
)

// NewServer returns a new web server
func NewServer(c Container) *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/html", htmlHandler(c))
	r.HandleFunc("/txt", txtHandler(c))

	port := 8080
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &srv
}

// htmlHandler returns a handler for serving a preview of our HTML message
func htmlHandler(c Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := executeMessageTemplateWithFallback("html", models.PreviewMessage(), c)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	}
}

// txtHandler returns a handler for serving a preview of our plain text message
func txtHandler(c Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := executeMessageTemplateWithFallback("txt", models.PreviewMessage(), c)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(content)
	}
}

// executeMessageTemplateWithFallback attempts to parse the message template that pertains to the provided type
// but falls back to example if not exists
func executeMessageTemplateWithFallback(messageTypeSuffix string, data interface{}, t TemplateInjector) ([]byte, error) {
	var buf bytes.Buffer

	if err := t.Template().ExecuteTemplate(&buf, fmt.Sprintf("message_%s", messageTypeSuffix), data); err != nil {
		// fallback to example template
		if err := t.Template().ExecuteTemplate(&buf, fmt.Sprintf("example_message_%s", messageTypeSuffix), data); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	return buf.Bytes(), nil
}
