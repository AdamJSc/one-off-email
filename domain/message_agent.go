package domain

import (
	"bytes"
	"errors"
	"fmt"
	htmlTpl "html/template"
	txtTpl "text/template"
)

type MessageAgent struct {
	tplClient TemplateClient
}

type TemplateClient interface {
	ParseHTMLTemplateFromFile(path string) (*htmlTpl.Template, error)
	ParseTxtTemplateFromFile(path string) (*txtTpl.Template, error)
}

func NewMessageAgent(tplClient TemplateClient) (*MessageAgent, error) {
	if tplClient == nil {
		return nil, errors.New("template client is nil")
	}

	m := &MessageAgent{
		tplClient: tplClient,
	}

	return m, nil
}

type GenerateMessageParams struct {
	TxtTemplatePath  string
	HTMLTemplatePath string
	Data             any
}

func (m *MessageAgent) Generate(p GenerateMessageParams) (*Message, error) {
	if p.TxtTemplatePath == "" {
		return nil, errors.New("txt template path is empty")
	}

	tbuf := &bytes.Buffer{}
	hbuf := &bytes.Buffer{}

	ttpl, err := m.tplClient.ParseTxtTemplateFromFile(p.TxtTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("cannot parse txt template: %w", err)
	}

	if err := ttpl.ExecuteTemplate(tbuf, "", p.Data); err != nil {
		return nil, fmt.Errorf("cannot execute txt template: %w", err)
	}

	if p.HTMLTemplatePath != "" {
		htpl, err := m.tplClient.ParseHTMLTemplateFromFile(p.HTMLTemplatePath)
		if err != nil {
			return nil, fmt.Errorf("cannot parse html template: %w", err)
		}

		if err := htpl.ExecuteTemplate(hbuf, "", p.Data); err != nil {
			return nil, fmt.Errorf("cannot execute html template: %w", err)
		}
	}

	msg := &Message{
		txtContent:  tbuf.String(),
		htmlContent: hbuf.String(),
	}

	return msg, nil
}
