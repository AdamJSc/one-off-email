package models

// Identity defines the recipient or sender identity of a Message
type Identity struct {
	Name  string
	Email string
}

// Message defines our message parameters
type Message struct {
	From   string
	To     string
	Advice string
}

// PreviewMessage returns a message to be used for previewing our templates
func PreviewMessage() *Message {
	return &Message{
		From:   "Eddie H",
		To:     "Jase T",
		Advice: "The more games you win or draw, the more points you'll receive.",
	}
}

// Email defines our email parameters
type Email struct {
	Sender    Identity
	Recipient Identity
	Message   Message
}
