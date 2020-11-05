package domain

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"one-off-email/app"
	"one-off-email/models"
)

// EmailAgentInjector defines the required behaviours for our EmailAgent
type EmailAgentInjector interface {
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

// ParseMessageTemplate attempts to parse the message template that pertains to the provided type
func (e *EmailAgent) ParseMessageTemplate(messageTypeSuffix string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if err := e.Template().ExecuteTemplate(&buf, fmt.Sprintf("message_%s", messageTypeSuffix), data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ParseMessageTemplateWithFallback attempts to parse the message template that pertains to the provided type
// but falls back to example if not exists
func (e *EmailAgent) ParseMessageTemplateWithFallback(messageTypeSuffix string, data interface{}) ([]byte, error) {
	content, err := e.ParseMessageTemplate(messageTypeSuffix, data)
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
