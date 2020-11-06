package models

import "fmt"

// Identity defines the recipient or sender identity of a Message
type Identity struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

func (i *Identity) String() string {
	if i.Email == "" {
		return ""
	}
	if i.Name == "" {
		return i.Email
	}
	return fmt.Sprintf("%s <%s>", i.Name, i.Email)
}

// RecipientList defines a list of recipient identities
type RecipientList []Identity

// Message defines our message parameters
type Message struct {
	From string
	To   string
}

// PreviewMessage returns a message to be used for previewing our templates
func PreviewMessage(from string) *Message {
	return &Message{
		From: from,
		To:   "Jase T",
	}
}

// Email defines our email parameters
type Email struct {
	Sender    Identity
	Recipient Identity
	ReplyTo   Identity
	BCC       Identity
	Subject   string
	Message   Message
}
