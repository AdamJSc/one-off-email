package main

import (
	"html/template"
	"log"
	"one-off-email/app"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	deps := dependencies{
		config:   app.MustParseConfig("data/.env"),
		template: app.MustParseTemplate("data/templates"),
	}

	deps.Template().ExecuteTemplate(log.Writer(), "message_html", nil)
}

// dependencies implements app.Dependencies
type dependencies struct {
	config   *app.Config
	template *template.Template
}

// Config implements app.Dependencies.Config()
func (d *dependencies) Config() *app.Config { return d.config }

// Config implements app.Dependencies.Template()
func (d *dependencies) Template() *template.Template { return d.template }
