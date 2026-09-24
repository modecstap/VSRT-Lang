package mail

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"mime"
	"net"
	"net/smtp"
	"strings"
)

type Sender struct {
	cfg Config
}

func NewSender(cfg Config) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) Send(to, link string) error {
	if err := s.cfg.Ready(); err != nil {
		return err
	}

	msg, err := RenderPasswordReset(link)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(s.cfg.Host, s.cfg.Port)
	encodedSubject := mime.QEncoding.Encode("utf-8", msg.Subject)
	payload := strings.Join([]string{
		"From: " + s.cfg.From,
		"To: " + to,
		"Subject: " + encodedSubject,
		"MIME-Version: 1.0",
		"Content-Type: " + msg.ContentType + "; charset=utf-8",
		"",
		msg.Body,
	}, "\r\n")

	if err := s.deliver(addr, to, []byte(payload)); err != nil {
		slog.Error("smtp send failed", "to", to, "error", err)
		return err
	}
	return nil
}

func (s *Sender) deliver(addr, to string, payload []byte) error {
	if s.cfg.Port == "465" {
		return s.deliverTLS(addr, to, payload)
	}
	return s.deliverStartTLS(addr, to, payload)
}

func (s *Sender) deliverTLS(addr, to string, payload []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: s.cfg.Host})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return err
	}
	defer client.Close()
	return s.transmit(client, to, payload)
}

func (s *Sender) deliverStartTLS(addr, to string, payload []byte) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: s.cfg.Host}); err != nil {
			return err
		}
	}
	return s.transmit(client, to, payload)
}

func (s *Sender) transmit(client *smtp.Client, to string, payload []byte) error {
	if s.cfg.User != "" {
		auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Pass, s.cfg.Host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}
	if err := client.Mail(s.cfg.From); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(payload); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("smtp quit: %w", err)
	}
	return nil
}
