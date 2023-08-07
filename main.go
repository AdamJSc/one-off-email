package main

import (
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

	// TODO: parse config
	c := &Config{}

	tplClient := tpl.NewClient()
	mgClient := mustMailgunClient(mailgun.NewClient())
	msgAgent := mustMessageAgent(domain.NewMessageAgent(tplClient))
	emailAgent := mustEmailAgent(domain.NewEmailAgent(mgClient))

	j := &job{
		txtTemplatePath:  txtTemplatePath,
		htmlTemplatePath: htmlTemplatePath,
		config:           c,
		msgAgent:         msgAgent,
		emailAgent:       emailAgent,
	}

	emails := j.mustParseEmails()

	if send {
		j.mustSendEmails(emails)
	} else {
		j.mustServeEmailsPreview(emails)
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

type job struct {
	txtTemplatePath  string
	htmlTemplatePath string
	config           *Config
	msgAgent         *domain.MessageAgent
	emailAgent       *domain.EmailAgent
}

func (j *job) mustParseEmails() []*domain.Email {
	msgConfigs := j.config.Messages
	emails := make([]*domain.Email, len(msgConfigs))

	for idx, msgConfig := range msgConfigs {
		p := domain.GenerateMessageParams{
			TxtTemplatePath:  j.txtTemplatePath,
			HTMLTemplatePath: j.htmlTemplatePath,
			Data:             msgConfig,
		}

		msg, err := j.msgAgent.Generate(p)
		if err != nil {
			log.Fatalf("cannot generate message index %d: %s", idx, err.Error())
		}

		// TODO: populate email field values once config has been parsed
		email, err := domain.NewEmail("", "", "", msg)
		if err != nil {
			log.Fatalf("cannot create email: message index %d: %s", idx, err.Error())
		}

		emails[idx] = email
	}

	return emails
}

func (j *job) mustServeEmailsPreview(emails []*domain.Email) {
	// TODO: implement must serve emails preview
}

func (j *job) mustSendEmails(emails []*domain.Email) {
	// TODO: implement must send emails
}

type Config struct {
	Messages []MessageConfig
}

type MessageConfig struct{}
