package whatsapp

import (
	"context"
	"github.com/21mebrat/lost-found-platform/internal/config"
)

type Request struct {
	Recipient string
	Message   string
}

type MetaProvider struct {
	cfg config.WhatsAppConfig
}

func NewMetaProvider(cfg config.WhatsAppConfig) (*MetaProvider, error) {
	return &MetaProvider{cfg: cfg}, nil
}

func (p *MetaProvider) Send(ctx context.Context, req Request) error {
	// Skeleton implementation
	return nil
}
