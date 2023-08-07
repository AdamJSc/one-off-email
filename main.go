package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/adamjsc/emailmerge/domain"
	"github.com/adamjsc/emailmerge/mailgun"
	"github.com/adamjsc/emailmerge/tpl"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var (
		templatesPath string
		send          bool
	)

	flag.StringVar(&templatesPath, "templates", "", "Specify path to templates folder (required)")
	flag.BoolVar(&send, "send", false, "If true, sends emails via mailgun (otherwise serves web preview of emails)")
	flag.Parse()

	if strings.Trim(templatesPath, " ") == "" {
		log.Fatal("-templates flag is required")
	}

	//configPath := filepath.Join(templatesPath, "config.json")
	txtTemplatePath := filepath.Join(templatesPath, "email.gotpl")
	htmlTemplatePath := filepath.Join(templatesPath, "email.gohtml")

	// TODO: parse config from configPath
	c := &config{}

	tplClient := tpl.NewClient()
	mgClient := mustMailgunClient(mailgun.NewClient())
	msgAgent := mustMessageAgent(domain.NewMessageAgent(tplClient))
	emailAgent := mustEmailAgent(domain.NewEmailAgent(mgClient))

	useCase, err := domain.NewEmailsUseCase(msgAgent, emailAgent, txtTemplatePath, htmlTemplatePath)
	if err != nil {
		log.Fatalf("cannot create emails use case: %w", err)
	}

	ctx := context.Background()

	if send {
		count, err := useCase.Issue(ctx, c.EmailConfigs)
		if err != nil {
			log.Fatalf("cannot issue emails: %s", err.Error())
		}

		log.Printf("%d emails issued successfully!", count)
	} else {
		if _, err := useCase.Preview(ctx, c.EmailConfigs); err != nil {
			log.Fatalf("cannot preview emails: %s", err.Error())
		}

		log.Println("finished!")
	}
}

func mustMessageAgent(agent *domain.MessageAgent, err error) *domain.MessageAgent {
	if err != nil {
		log.Fatal(err)
	}

	return agent
}

func mustEmailAgent(agent *domain.EmailAgent, err error) *domain.EmailAgent {
	if err != nil {
		log.Fatal(err)
	}

	return agent
}

func mustMailgunClient(mg *mailgun.Client, err error) *mailgun.Client {
	if err != nil {
		log.Fatal(err)
	}

	return mg
}

type config struct {
	EmailConfigs []domain.EmailConfig
}
