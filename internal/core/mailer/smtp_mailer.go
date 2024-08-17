package mailer

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/smtp"
	"strconv"
)

var _ Sender = (*SMTPSender)(nil)

type SMTPConfig struct {
	Host     string `json:"host,omitempty"     conf:"default:localhost"`
	Port     int    `json:"port,omitempty"     conf:"default:1025"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	From     string `json:"from,omitempty"     conf:"Gotll Inc."`
}

func (m *SMTPConfig) Ready() bool {
	return m.Host != "" && m.Port != 0 && m.Username != "" && m.Password != "" && m.From != ""
}

func (m *SMTPConfig) server() string {
	return m.Host + ":" + strconv.Itoa(m.Port)
}

type SMTPSender struct {
	cfg SMTPConfig
}

func NewSMTPSender(cfg SMTPConfig) *SMTPSender {
	return &SMTPSender{cfg: cfg}
}

func (m *SMTPSender) Send(msg Message) error {
	server := m.cfg.server()

	header := map[string]string{
		"From":                      msg.From,
		"Subject":                   mime.QEncoding.Encode("UTF-8", msg.Subject),
		"MIME-Version":              "1.0",
		"Content-Type":              "text/plain; charset=\"utf-8\"",
		"Content-Transfer-Encoding": "base64",
	}

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(msg.Body))

	return smtp.SendMail(
		server,
		smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host),
		m.cfg.From,
		[]string{msg.To},
		[]byte(message),
	)
}
