package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResendProvider struct {
	config *Config
	client *http.Client
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type resendEmailResponse struct {
	ID string `json:"id"`
}

type resendErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Name       string `json:"name"`
}

func NewResendProvider(cfg *Config) (Provider, error) {
	if cfg.Password == "" {
		return nil, fmt.Errorf("resend API key is required")
	}

	return &ResendProvider{
		config: cfg,
		client: &http.Client{},
	}, nil
}

func (p *ResendProvider) Send(ctx context.Context, req *EmailRequest) (*EmailResponse, error) {
	body := req.Body
	if req.HTMLBody != "" {
		body = req.HTMLBody
	}

	resendReq := resendEmailRequest{
		From:    req.From,
		To:      []string{req.To},
		Subject: req.Subject,
		HTML:    body,
	}

	jsonData, err := json.Marshal(resendReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.resend.com/emails",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.config.Password)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		var errResp resendErrorResponse
		if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
			return &EmailResponse{
				MessageID: req.MessageID,
				Status:    "Failed",
				Error:     fmt.Sprintf("resend API error %d: %s", resp.StatusCode, string(bodyBytes)),
			}, fmt.Errorf("resend API error %d", resp.StatusCode)
		}
		return &EmailResponse{
			MessageID: req.MessageID,
			Status:    "Failed",
			Error:     errResp.Message,
		}, fmt.Errorf("resend error: %s", errResp.Message)
	}

	var emailResp resendEmailResponse
	if err := json.Unmarshal(bodyBytes, &emailResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &EmailResponse{
		MessageID: emailResp.ID,
		Status:    "Sent",
		Error:     "",
	}, nil
}

func (p *ResendProvider) VerifyCredentials(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.resend.com/domains",
		nil,
	)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.config.Password)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("invalid Resend API key")
	}

	return nil
}

func (p *ResendProvider) GetProviderName() string {
	return "resend"
}
