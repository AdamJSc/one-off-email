package domain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"one-off-email/app"
	"one-off-email/models"

	"gopkg.in/yaml.v2"
)

// mailgunSuccessMessage defines a known good response message from the Mailgun SDK client
const mailgunSuccessMessage = "Queued. Thank you."

// EmailAgentInjector defines the required behaviours for our EmailAgent
type EmailAgentInjector interface {
	app.ConfigInjector
	app.TemplateInjector
	app.MailgunInjector
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
	cfg := e.Config()

	return &models.Email{
		Sender: models.Identity{
			Name:  cfg.SenderName,
			Email: cfg.SenderEmail,
		},
		Recipient: recipient,
		ReplyTo: models.Identity{
			Name:  cfg.ReplyToName,
			Email: cfg.ReplyToEmail,
		},
		BCC: models.Identity{
			Name:  cfg.BCCName,
			Email: cfg.BCCEmail,
		},
		Subject: cfg.EmailSubject,
		Message: models.Message{
			From: cfg.MessageSignOff,
			To:   recipient.Name,
		},
	}
}

// IssueEmail issues the provided email
func (e *EmailAgent) IssueEmail(ctx context.Context, email *models.Email) error {
	mg := e.Mailgun()

	bodyPlain, err := parseMessageTemplate("txt", email.Message, e.Template())
	if err != nil {
		// must have plain text body
		return err
	}

	// new mailgun message
	mgMsg := mg.NewMessage(
		email.Sender.String(),
		email.Subject,
		string(bodyPlain),
		email.Recipient.String(),
	)

	bodyHTML, err := parseMessageTemplate("html", email.Message, e.Template())
	if err == nil {
		// add html body if we have one
		mgMsg.SetHtml(string(bodyHTML))
	}

	mgMsg.SetTracking(true)

	replyTo := email.ReplyTo.String()
	if replyTo != "" {
		mgMsg.SetReplyTo(replyTo)
	}

	bcc := email.BCC.String()
	if bcc != "" {
		mgMsg.AddBCC(bcc)
	}

	result, id, err := mg.Send(ctx, mgMsg)
	if err != nil {
		return err
	}
	if result != mailgunSuccessMessage {
		return fmt.Errorf("send message result: %s", result)
	}
	if id == "" {
		return errors.New("send id empty")
	}

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
