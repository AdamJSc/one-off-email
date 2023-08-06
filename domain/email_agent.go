package domain

import (
	"context"
	"errors"
	"fmt"
)

type EmailAgent struct {
	client EmailClient
}

type EmailClient interface {
	SendEmail(ctx context.Context, email *Email) (string, error)
}

func NewEmailAgent(client EmailClient) (*EmailAgent, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	e := &EmailAgent{
		client: client,
	}

	return e, nil
}

func (e *EmailAgent) Send(ctx context.Context, email *Email) (string, error) {
	// TODO: monitor cancelled context

	res, err := e.client.SendEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("cannot send email: %w", err)
	}

	return res, nil
}
