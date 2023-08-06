package handlers

import (
	"fmt"
	"log"
	"net/http"
	"one-off-email/app"
	"one-off-email/domain"
	"one-off-email/models"
	"time"

	"github.com/gorilla/mux"
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
	agent := domain.EmailAgentLegacy{EmailAgentInjector: c}
	config := c.Config()
	subject := config.EmailSubject
	from := config.MessageSignOff

	return func(w http.ResponseWriter, r *http.Request) {
		heading := fmt.Sprintf(`<p>[Subject]<br />%s</p><p>[Body]<br /></p>`, subject)

		content, err := agent.ParseMessageTemplateWithFallback("html", models.PreviewMessage(from))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(heading))
		w.Write(content)
	}
}

// txtHandler returns a handler for serving a preview of our plain text message
func txtHandler(c app.Container) http.HandlerFunc {
	agent := domain.EmailAgentLegacy{EmailAgentInjector: c}
	config := c.Config()
	subject := config.EmailSubject
	from := config.MessageSignOff

	return func(w http.ResponseWriter, r *http.Request) {
		heading := fmt.Sprintf(`[Subject]
%s

[Body]
`, subject)

		content, err := agent.ParseMessageTemplateWithFallback("txt", models.PreviewMessage(from))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(heading))
		w.Write(content)
	}
}
