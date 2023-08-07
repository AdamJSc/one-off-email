package domain

import (
	"context"
	"errors"
	"fmt"
)

type EmailsUseCase struct {
	MessageAgent     *MessageAgent
	EmailAgent       *EmailAgent
	TxtTemplatePath  string
	HTMLTemplatePath string
}

func NewEmailsUseCase(msgAgent *MessageAgent, emailAgent *EmailAgent, txtTemplatePath, htmlTemplatePath string) (*EmailsUseCase, error) {
	if msgAgent == nil {
		return nil, errors.New("message agent is nil")
	}
	if emailAgent == nil {
		return nil, errors.New("email agent is nil")
	}
	if txtTemplatePath == "" {
		return nil, errors.New("txt template path is empty")
	}

	e := &EmailsUseCase{
		MessageAgent:     msgAgent,
		EmailAgent:       emailAgent,
		TxtTemplatePath:  txtTemplatePath,
		HTMLTemplatePath: htmlTemplatePath,
	}

	return e, nil
}

type EmailConfig struct{}

func (e *EmailsUseCase) Issue(ctx context.Context, configs []EmailConfig) (int, error) {
	emails, err := generateEmails(e.MessageAgent, e.TxtTemplatePath, e.HTMLTemplatePath, configs)
	if err != nil {
		return 0, fmt.Errorf("cannot genrate emails: %w", err)
	}

	// TODO: issue emails via email agent

	return len(emails), nil
}

func generateEmails(msgAgent *MessageAgent, txtTemplatePath, htmlTemplatePath string, configs []EmailConfig) ([]*Email, error) {
	emails := make([]*Email, len(configs))

	for idx, config := range configs {
		p := GenerateMessageParams{
			TxtTemplatePath:  txtTemplatePath,
			HTMLTemplatePath: htmlTemplatePath,
			Data:             config,
		}

		msg, err := msgAgent.Generate(p)
		if err != nil {
			return nil, fmt.Errorf("cannot generate message: config index %d: %s", idx, err.Error())
		}

		// TODO: populate email field values once config has been parsed
		email, err := NewEmail("", "", "", msg)
		if err != nil {
			return nil, fmt.Errorf("cannot create email: config index %d: %s", idx, err.Error())
		}

		emails[idx] = email
	}

	return emails, nil
}

func (e *EmailsUseCase) Preview(ctx context.Context, configs []EmailConfig) (int, error) {
	emails, err := generateEmails(e.MessageAgent, e.TxtTemplatePath, e.HTMLTemplatePath, configs)
	if err != nil {
		return 0, fmt.Errorf("cannot genrate emails: %w", err)
	}

	// TODO: render email preview via web server

	return len(emails), nil
}
