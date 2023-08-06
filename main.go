package main

import (
	"flag"
	"log"
	"one-off-email/domain"
	"one-off-email/mailgun"
	"one-off-email/tpl"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var (
		configPath       string
		txtTemplatePath  string
		htmlTemplatePath string
		send             bool
	)

	flag.StringVar(&configPath, "config", "", "Specify path to config file (required)")
	flag.StringVar(&txtTemplatePath, "txt", "", "Specify path to txt template file (required)")
	flag.StringVar(&htmlTemplatePath, "html", "", "Specify path to HTML template file (optional)")
	flag.BoolVar(&send, "send", false, "If true, sends emails via mailgun (otherwise serves web preview of emails)")
	flag.Parse()

	if strings.Trim(configPath, " ") == "" {
		log.Fatal("-config flag is required")
	}
	if strings.Trim(txtTemplatePath, " ") == "" {
		log.Fatal("-txt flag is required")
	}

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
