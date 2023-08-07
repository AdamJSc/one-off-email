package mailgun

import (
	"context"

	"github.com/adamjsc/emailmerge/domain"
)

type Client struct{}

func NewClient() (*Client, error) {
	// TODO: implement new mailgun client
	return &Client{}, nil
}

func (c *Client) SendEmail(ctx context.Context, email *domain.Email) (string, error) {
	// TODO: implement mailgun send email method
	return "", nil
}
