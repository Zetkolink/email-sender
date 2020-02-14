package helpers

import (
	"email-sender/models"
	"net/smtp"
	"strings"
)

type SmtpHandler struct {
	auth smtp.Auth
}

func InitSmtp(cfg Config) *SmtpHandler {
	auth := smtp.PlainAuth(
		cfg.Smtp.Identity,
		cfg.Smtp.Username,
		cfg.Smtp.Password,
		cfg.Smtp.Hostname,
	)

	return &SmtpHandler{
		auth: auth,
	}
}

func (s SmtpHandler) buildBody(msg models.MessageRequest) []byte {
	body := "From:" + msg.Sender + "\n" + "To:" + strings.Join(msg.To, ",") + "\n"
	if msg.Subject != nil {
		body += "Subject:" + *msg.Subject + "\n"
	}
	body += msg.Message

	return []byte(body)
}

func (s SmtpHandler) SendMail(addr string, msg models.MessageRequest) error {
	err := smtp.SendMail(
		addr,
		s.auth,
		msg.Sender,
		msg.To,
		s.buildBody(msg),
	)
	if err != nil {
		return err
	}

	return nil
}
