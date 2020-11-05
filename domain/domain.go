package domain

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"one-off-email/app"
	"one-off-email/models"
)

// EmailAgentInjector defines the required behaviours for our EmailAgent
type EmailAgentInjector interface {
	app.ConfigInjector
	app.TemplateInjector
}

// EmailAgent defines all of our email-related operations
type EmailAgent struct {
	EmailAgentInjector
}

// ParseRecipientsFromFile parses a recipient list from the provided file path
func (e *EmailAgent) ParseRecipientsFromFile(path string) (models.RecipientList, error) {
	var fileContents struct {
		Recipients models.RecipientList `yaml:"recipients"`
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(bytes, &fileContents); err != nil {
		return nil, err
	}

	return fileContents.Recipients, nil
}

// ParseMessageTemplateWithFallback attempts to parse the message template that pertains to the provided type
// but falls back to example if not exists
func (e *EmailAgent) ParseMessageTemplateWithFallback(messageTypeSuffix string, data interface{}) ([]byte, error) {
	content, err := parseMessageTemplate(messageTypeSuffix, data, e.Template())
	if err == nil {
		return content, nil
	}

	// otherwise fallback to example template
	var buf bytes.Buffer
	if err := e.Template().ExecuteTemplate(&buf, fmt.Sprintf("example_message_%s", messageTypeSuffix), data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateEmail generates an email object from the provided recipient
func (e *EmailAgent) GenerateEmail(recipient models.Identity) *models.Email {
	return &models.Email{
		Sender: models.Identity{
			Name:  e.Config().SenderName,
			Email: e.Config().SenderEmail,
		},
		Recipient: recipient,
		Message: models.Message{
			From: e.Config().MessageSignOff,
			To:   recipient.Name,
		},
	}
}

// IssueEmail issues the provided email
func (e *EmailAgent) IssueEmail(email *models.Email) error {
	// TODO - implement me
	return nil
}

// parseMessageTemplate attempts to parse the message template that pertains to the provided type
func parseMessageTemplate(messageTypeSuffix string, data interface{}, t *template.Template) ([]byte, error) {
	var buf bytes.Buffer

	if err := t.ExecuteTemplate(&buf, fmt.Sprintf("message_%s", messageTypeSuffix), data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
