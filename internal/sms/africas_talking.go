package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/21mebrat/lost-found-platform/internal/config"
)

type Request struct {
	Recipient string
	Message   string
}

type AfricasTalkingProvider struct {
	apiKey   string
	username string
	senderID string
	baseURL  string
	client   *http.Client
}

func NewAfricasTalkingProvider(
	cfg config.SMSConfig,
	timeout time.Duration,
) (*AfricasTalkingProvider, error) {

	// if strings.TrimSpace(cfg.APIKey) == "" {
	// 	return nil, errors.New("SMS API key is required")
	// }

	// if strings.TrimSpace(cfg.Username) == "" {
	// 	return nil, errors.New("SMS username is required")
	// }

	if strings.TrimSpace(cfg.SenderID) == "" {
		return nil, errors.New("SMS sender ID is required")
	}

	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, errors.New("SMS base URL is required")
	}

	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &AfricasTalkingProvider{
		apiKey:   cfg.APIKey,
		username: cfg.Username,
		senderID: cfg.SenderID,
		baseURL:  strings.TrimRight(cfg.BaseURL, "/"),
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

type sendRequest struct {
	Username string `json:"username"`
	To       string `json:"to"`
	Message  string `json:"message"`
	From     string `json:"from"`
}

type sendResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (p *AfricasTalkingProvider) Send(
	ctx context.Context,
	req Request,
) error {

	if err := validateRequest(req); err != nil {
		return err
	}

	payload := sendRequest{
		Username: p.username,
		To:       req.Recipient,
		Message:  req.Message,
		From:     p.senderID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf(
			"marshal SMS request: %w",
			err,
		)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.baseURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf(
			"create SMS request: %w",
			err,
		)
	}

	httpReq.Header.Set(
		"Content-Type",
		"application/json",
	)

	httpReq.Header.Set(
		"Accept",
		"application/json",
	)

	httpReq.Header.Set(
		"apiKey",
		p.apiKey,
	)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf(
			"send SMS request: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {

		return fmt.Errorf(
			"SMS provider returned HTTP status %d",
			resp.StatusCode,
		)
	}

	var result sendResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf(
			"decode SMS provider response: %w",
			err,
		)
	}

	if !result.Success {
		return fmt.Errorf(
			"SMS provider rejected request: %s",
			result.Message,
		)
	}

	return nil
}

func validateRequest(req Request) error {
	if strings.TrimSpace(req.Recipient) == "" {
		return errors.New("SMS recipient is required")
	}

	if strings.TrimSpace(req.Message) == "" {
		return errors.New("SMS message is required")
	}

	return nil
}
