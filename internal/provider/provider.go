package provider

import (
	"context"
	"fmt"
)

type EmailRequest struct {
	MessageID string
	From      string
	To        string
	Subject   string
	Body      string
	HTMLBody  string
}

type EmailResponse struct {
	MessageID string
	Status    string
	Error     string
}

type Provider interface {
	Send(ctx context.Context, req *EmailRequest) (*EmailResponse, error)
	VerifyCredentials(ctx context.Context) error
	GetProviderName() string
}

type Config struct {
	Provider    string
	Host        string
	Port        int
	Username    string
	Password    string
	Timeout     int
	SenderEmail string
}

func NewProvider(cfg *Config) (Provider, error) {
	switch cfg.Provider {
	case "native-smtp":
		return NewNativeSMTPProvider(cfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}
