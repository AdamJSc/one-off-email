package main

import (
	"flag"
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

	send := flag.Bool("send", false, "include this flag to physically issue emails")
	flag.Parse()

	if !(*send) {
		// run in preview mode
		srv := app.NewServer(&deps)
		log.Printf("listening on %s...\n", srv.Addr)
		log.Fatal(srv.ListenAndServe())
	}

	// send emails
	log.Println("sending emails...")
}

// dependencies implements app.Container
type dependencies struct {
	config   *app.Config
	template *template.Template
}

// Config implements app.Container.Config()
func (d *dependencies) Config() *app.Config { return d.config }

// Config implements app.Container.Template()
func (d *dependencies) Template() *template.Template { return d.template }
