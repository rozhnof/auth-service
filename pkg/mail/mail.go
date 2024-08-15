package mail

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/pkg/errors"
)

type MailSender struct {
	senderEmail    string
	senderPassword string
	host           string
}

func NewMailSender(senderEmail string, senderPassword string, host string) *MailSender {
	return &MailSender{
		senderEmail:    senderEmail,
		senderPassword: senderPassword,
		host:           host,
	}
}

func (s *MailSender) SendMessage(sender string, recipient string, subject string, body string) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", sender, recipient, subject, body)

	host := strings.Split(s.host, ":")
	if len(host) != 2 {
		return errors.Wrap(ErrInvalidMailHost, s.host)
	}
	port := host[0]

	err := smtp.SendMail(s.host,
		smtp.PlainAuth("", s.senderEmail, s.senderPassword, port),
		s.senderEmail, []string{recipient}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}
