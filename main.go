package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"one-off-email/domain"
	"one-off-email/mailgun"
	"one-off-email/tpl"
	"time"
)

func main() {
	var (
		txtTemplatePath  string
		htmlTemplatePath string
	)

	flag.StringVar(&txtTemplatePath, "txt", "", "Specify path to Txt template file (required)")
	flag.StringVar(&htmlTemplatePath, "html", "", "Specify path to HTML template file (optional)")
	// TODO: parse config file path
	flag.Parse()

	if txtTemplatePath == "" {
		log.Fatalf("-txt flag is required")
	}
	// TODO: ensure config file path is not empty

	// TODO: parse config
	c := &config{}

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

	j.do()
}

type job struct {
	txtTemplatePath  string
	htmlTemplatePath string
	config           *config
	msgAgent         *domain.MessageAgent
	emailAgent       *domain.EmailAgent
}

func (j *job) do() {
	for idx, msgConfig := range j.config.messages {
		res, err := j.processMessageConfig(msgConfig)
		if err != nil {
			log.Fatalf("message config index %d: %s", idx, err.Error())
		}
		log.Printf("message send sucessful: message config index %d: sent id %q", idx, res)
	}
}

func (j *job) processMessageConfig(msgConfig messageConfig) (string, error) {
	msg, err := j.msgAgent.Generate(domain.GenerateMessageParams{
		TxtTemplatePath:  j.txtTemplatePath,
		HTMLTemplatePath: j.htmlTemplatePath,
		Data:             msgConfig,
	})
	if err != nil {
		return "", fmt.Errorf("cannot generate message: %w", err)
	}

	// TODO: populate email field values
	email, err := domain.NewEmail("", "", "", msg)
	if err != nil {
		return "", fmt.Errorf("cannot create email: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := j.emailAgent.Send(ctx, email)
	if err != nil {
		return "", fmt.Errorf("cannot send email: %w", err)
	}

	return res, nil
}

type config struct {
	messages []messageConfig
}

type messageConfig struct{}

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
