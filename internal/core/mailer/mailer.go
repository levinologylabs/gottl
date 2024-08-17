// Package mailer provides a simple interface for sending emails.
package mailer

type Message struct {
	To      string
	From    string
	Subject string
	Body    string
}

type Sender interface {
	Send(mail Message) error
}
