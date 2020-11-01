package main

import (
	"flag"
	"html/template"
	"log"
	"one-off-email/app"
	"one-off-email/domain"
	"one-off-email/models"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	deps := dependencies{
		config:   app.MustParseConfig("data/.env"),
		template: app.MustParseTemplate("data/templates"),
	}

	send := flag.Bool("send", false, "include this flag to physically issue emails")
	flag.Parse()

	if *send {
		sendEmails()
		return
	}

	// run in preview mode
	srv := app.NewServer(&deps)
	log.Printf("listening on %s...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func sendEmails() {
	var (
		list models.RecipientList
		err  error
	)

	agent := domain.EmailAgent{}

	// try to parse recipients
	list, err = agent.ParseRecipientsFromFile("data/recipients.yml")
	if err != nil {
		// parse example recipients
		list, err = agent.ParseRecipientsFromFile("data/recipients.example.yml")
		if err != nil {
			log.Fatal(err)
		}
	}

	// do stuff with recipients
	log.Println(list)
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
