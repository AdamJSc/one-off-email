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
		var buf bytes.Buffer
		if err := c.Template().ExecuteTemplate(&buf, "message_html", models.PreviewMessage()); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write(buf.Bytes())
	}
}

// txtHandler returns a handler for serving a preview of our plain text message
func txtHandler(c Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if err := c.Template().ExecuteTemplate(&buf, "message_txt", models.PreviewMessage()); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(buf.Bytes())
	}
}
