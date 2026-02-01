package sender

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	"delayed_notifier/internal/models"
)

// EmailSender отправляет email от имени сервиса
type EmailSender struct {
	Host     string
	Port     int
	Username string // email от которого отправляем
	Password string // пароль приложения
}

func NewEmailSender(host string, port int, username, password string) *EmailSender {
	return &EmailSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (e *EmailSender) Send(n models.Notification) error {
	addr := fmt.Sprintf("%s:%d", e.Host, e.Port)
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)

	msg := []byte(
		fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n",
			e.Username,
			n.Recipient,
			n.Subject,
			n.Message,
		),
	)

	// STARTTLS (587)
	conn, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Quit(); err != nil {
			log.Printf("failed to quit SMTP connection: %v", err)
		}
	}()

	if ok, _ := conn.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: e.Host}
		if err := conn.StartTLS(tlsConfig); err != nil {
			return err
		}
	}

	if err := conn.Auth(auth); err != nil {
		return err
	}

	if err := conn.Mail(e.Username); err != nil {
		return err
	}

	if err := conn.Rcpt(n.Recipient); err != nil {
		return err
	}

	w, err := conn.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
}
