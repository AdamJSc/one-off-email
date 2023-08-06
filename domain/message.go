package domain

import (
	"errors"
	"strings"
)

type Message struct {
	txtContent  string
	htmlContent string
}

func (m *Message) validate() error {
	if strings.Trim(m.txtContent, " ") == "" {
		return errors.New("txt content is empty")
	}

	return nil
}
