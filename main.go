package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"log"
	"one-off-email/app"
	"one-off-email/domain"
	"one-off-email/handlers"
	"one-off-email/models"
	"os"
	"sync"
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
		// not running in preview mode
		sendEmails(&deps)
		return
	}

	// run in preview mode
	srv := handlers.NewServer(&deps)
	log.Printf("listening on %s...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// sendEmails contains our entry-point logic for issuing our emails
func sendEmails(c app.Container) {
	var (
		recipients models.RecipientList
		err        error
	)

	agent := domain.EmailAgent{
		EmailAgentInjector: c,
	}

	// try to parse recipients
	recipients, err = agent.ParseRecipientsFromFile("data/recipients.yml")
	if err != nil {
		// parse example recipients
		recipients, err = agent.ParseRecipientsFromFile("data/recipients.example.yml")
		if err != nil {
			log.Fatal(err)
		}
	}

	// prompt user to continue
	log.Printf("send email to %d recipients? [y/N]", len(recipients))
	inp := bufio.NewScanner(os.Stdin)
	inp.Scan()
	if inp.Text() != "y" {
		log.Fatal("process aborted")
	}

	// do stuff with recipients
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // max concurrent processes
	errChan := make(chan error)

	for _, r := range recipients {
		wg.Add(1)
		sem <- struct{}{}
		go issueEmail(r, &agent, &wg, sem, errChan)
	}

	wg.Wait()
	close(sem)
	close(errChan)

	log.Printf("process complete... %d processed, %d failed\n", len(recipients), len(errChan))
	for err := range errChan {
		log.Println(err)
	}
}

// issueEmail issues an email to the provided recipient
func issueEmail(
	recipient models.Identity,
	agent *domain.EmailAgent,
	wg *sync.WaitGroup,
	sem chan struct{},
	errChan chan error,
) {
	log.Printf("issuing email to %s...", recipient.Email)

	var done = func() {
		wg.Done()
		<-sem
	}

	email := agent.GenerateEmail(recipient)

	if err := agent.IssueEmail(email); err != nil {
		errChan <- fmt.Errorf("error issuing to %s: %s", recipient.Email, err.Error())
		done()
		return
	}

	done()
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
