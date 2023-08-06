package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type Email struct {
	from    *mail.Address
	to      *mail.Address
	cc      *mail.Address
	bcc     *mail.Address
	replyTo *mail.Address
	subject string
	message *Message
}

func NewEmail(from, to, subject string, message *Message, opts ...EmailOpt) (*Email, error) {
	fromAddr, err := mail.ParseAddress(from)
	if err != nil {
		return nil, fmt.Errorf("invalid from address %q: %w", from, err)
	}
	toAddr, err := mail.ParseAddress(to)
	if err != nil {
		return nil, fmt.Errorf("invalid to address %q: %w", to, err)
	}
	subject = strings.Trim(subject, " ")
	if subject == "" {
		return nil, errors.New("subject is empty")
	}
	if message == nil {
		return nil, errors.New("message is nil")
	}
	if err := message.validate(); err != nil {
		return nil, fmt.Errorf("message invalid: %w", err)
	}

	e := &Email{
		from:    fromAddr,
		to:      toAddr,
		subject: subject,
		message: message,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}

	return e, nil
}

type EmailOpt func(e *Email) error

func SetCC(cc string) EmailOpt {
	return func(e *Email) error {
		if e == nil {
			return nil
		}
		addr, err := mail.ParseAddress(cc)
		if err != nil {
			return fmt.Errorf("invalid cc address: %w", err)
		}
		e.cc = addr
		return nil
	}
}

func SetBCC(bcc string) EmailOpt {
	return func(e *Email) error {
		if e == nil {
			return nil
		}
		addr, err := mail.ParseAddress(bcc)
		if err != nil {
			return fmt.Errorf("invalid bcc address: %w", err)
		}
		e.bcc = addr
		return nil
	}
}

func SetReplyTo(replyTo string) EmailOpt {
	return func(e *Email) error {
		if e == nil {
			return nil
		}
		addr, err := mail.ParseAddress(replyTo)
		if err != nil {
			return fmt.Errorf("invalid reply-to address: %w", err)
		}
		e.replyTo = addr
		return nil
	}
}
