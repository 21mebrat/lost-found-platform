package telegram

import (
	"context"
	"github.com/21mebrat/lost-found-platform/internal/config"
)

type Request struct {
	Recipient string
	Message   string
}

type TelegramProvider struct {
	cfg config.TelegramConfig
}

func NewTelegramProvider(cfg config.TelegramConfig) (*TelegramProvider, error) {
	return &TelegramProvider{cfg: cfg}, nil
}

func (p *TelegramProvider) Send(ctx context.Context, req Request) error {
	// Skeleton implementation
	return nil
}
