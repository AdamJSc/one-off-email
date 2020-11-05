package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"one-off-email/app"
	"one-off-email/domain"
	"one-off-email/models"
	"time"
)

// NewServer returns a new web server
func NewServer(c app.Container) *http.Server {
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
func htmlHandler(c app.Container) http.HandlerFunc {
	agent := domain.EmailAgent{EmailAgentInjector: c}

	return func(w http.ResponseWriter, r *http.Request) {
		content, err := agent.ParseMessageTemplateWithFallback("html", models.PreviewMessage(c.Config().MessageSignOff))
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
func txtHandler(c app.Container) http.HandlerFunc {
	agent := domain.EmailAgent{EmailAgentInjector: c}

	return func(w http.ResponseWriter, r *http.Request) {
		content, err := agent.ParseMessageTemplateWithFallback("txt", models.PreviewMessage(c.Config().MessageSignOff))
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
