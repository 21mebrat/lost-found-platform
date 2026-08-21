package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/21mebrat/lost-found-platform/internal/config"
)

type Request struct {
	Recipient string
	Subject   string
	Message   string
}

type SMTPProvider struct {
	cfg config.EmailConfig
}

func NewSMTPProvider(cfg config.EmailConfig) (*SMTPProvider, error) {
	if strings.TrimSpace(cfg.SMTPHost) == "" {
		return nil, fmt.Errorf("SMTP host is required")
	}
	if strings.TrimSpace(cfg.SMTPPort) == "" {
		return nil, fmt.Errorf("SMTP port is required")
	}
	if strings.TrimSpace(cfg.From) == "" {
		return nil, fmt.Errorf("SMTP sender email (From) is required")
	}
	return &SMTPProvider{cfg: cfg}, nil
}

func (p *SMTPProvider) Send(ctx context.Context, req Request) error {
	if strings.TrimSpace(req.Recipient) == "" {
		return fmt.Errorf("email recipient is required")
	}

	var auth smtp.Auth
	if p.cfg.SMTPUsername != "" || p.cfg.SMTPPassword != "" {
		auth = smtp.PlainAuth("", p.cfg.SMTPUsername, p.cfg.SMTPPassword, p.cfg.SMTPHost)
	}

	addr := fmt.Sprintf("%s:%s", p.cfg.SMTPHost, p.cfg.SMTPPort)

	// Format simple RFC 822 email message
	msg := fmt.Appendf(nil,
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		p.cfg.From,
		req.Recipient,
		req.Subject,
		req.Message,
	)

	err := smtp.SendMail(addr, auth, p.cfg.From, []string{req.Recipient}, msg)
	if err != nil {
		return fmt.Errorf("send SMTP email: %w", err)
	}

	return nil
}
