package domain

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"one-off-email/models"
)

// EmailAgent defines all of our email-related operations
type EmailAgent struct{}

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
