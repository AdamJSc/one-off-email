package models

// Identity defines the recipient or sender identity of a Message
type Identity struct {
	Name  string
	Email string
}

// Message defines our message parameters
type Message struct {
	Sender    Identity
	Recipient Identity
	Advice    string
}

// PreviewMessage returns a message to be used for previewing our templates
func PreviewMessage() *Message {
	return &Message{
		Recipient: Identity{
			Name:  "Jase",
			Email: "tindall_tindall_give_us_a_wave@example.net",
		},
		Sender: Identity{
			Name:  "Eddie",
			Email: "eddie_eddie_eddie_howe@example.net",
		},
		Advice: "The more games you win or draw, the more points you'll receive.",
	}
}

// Email defines our email parameters
type Email struct {
	From    Identity
	To      Identity
	Message Message
}
