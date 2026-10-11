package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

// SMTPConfig is the outbound mail server. Port 465 uses implicit TLS.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// Enabled reports whether a message can actually be sent.
func (c SMTPConfig) Enabled() bool {
	return c.Host != "" && c.Port > 0 && c.Username != "" && c.Password != "" && c.From != ""
}

// Sender delivers plain-text mail over implicit TLS (SMTPS).
type Sender struct {
	cfg SMTPConfig
}

// NewSender builds an SMTP sender. A disabled config still returns a sender
// whose Send method refuses to connect.
func NewSender(cfg SMTPConfig) *Sender {
	return &Sender{cfg: cfg}
}

// Enabled reports whether SMTP settings are complete.
func (s *Sender) Enabled() bool { return s.cfg.Enabled() }

// Send delivers one plain-text message. It does not retry.
func (s *Sender) Send(ctx context.Context, to, subject, text string) error {
	if !s.cfg.Enabled() {
		return fmt.Errorf("smtp is not configured")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	dialer := &tls.Dialer{Config: &tls.Config{ServerName: s.cfg.Host, MinVersion: tls.VersionTLS12}}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if err := client.Auth(smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(s.cfg.From); err != nil {
		return fmt.Errorf("smtp from: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	message := strings.Join([]string{
		"From: " + s.cfg.From,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		text,
	}, "\r\n")
	if _, err := writer.Write([]byte(message)); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}
	return client.Quit()
}
