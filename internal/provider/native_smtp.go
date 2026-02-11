package provider

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

type NativeSMTPProvider struct {
	host     string
	port     int
	username string
	password string
	timeout  time.Duration
}

func NewNativeSMTPProvider(cfg *Config) (Provider, error) {
	if cfg.Host == "" || cfg.Port == 0 {
		return nil, fmt.Errorf("host and port required")
	}
	timeout := time.Duration(cfg.Timeout) * time.Second
	if cfg.Timeout == 0 {
		timeout = 30 * time.Second
	}
	return &NativeSMTPProvider{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		timeout:  timeout,
	}, nil
}

func (p *NativeSMTPProvider) Send(ctx context.Context, req *EmailRequest) (*EmailResponse, error) {
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, p.host)
	if err != nil {
		return nil, fmt.Errorf("smtp client failed: %w", err)
	}
	defer client.Close()

	if err := client.Mail(req.From); err != nil {
		return nil, fmt.Errorf("mail from failed: %w", err)
	}

	if err := client.Rcpt(req.To); err != nil {
		return nil, fmt.Errorf("rcpt to failed: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return nil, fmt.Errorf("data failed: %w", err)
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", req.From, req.To, req.Subject, req.Body)
	if _, err := wc.Write([]byte(msg)); err != nil {
		wc.Close()
		return nil, fmt.Errorf("write failed: %w", err)
	}

	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("close failed: %w", err)
	}

	_ = client.Quit()

	return &EmailResponse{MessageID: req.MessageID, Status: "Sent"}, nil
}

func (p *NativeSMTPProvider) VerifyCredentials(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, p.host)
	if err != nil {
		return fmt.Errorf("smtp client failed: %w", err)
	}
	defer client.Close()

	return nil
}

func (p *NativeSMTPProvider) GetProviderName() string {
	return "native-smtp"
}
