package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/21mebrat/lost-found-platform/internal/config"
	"github.com/21mebrat/lost-found-platform/internal/email"
	"github.com/21mebrat/lost-found-platform/internal/sms"
	"github.com/21mebrat/lost-found-platform/internal/telegram"
	"github.com/21mebrat/lost-found-platform/internal/whatsapp"
)

type Factory struct {
	cfg config.NotificationConfig
}

func NewFactory(cfg config.NotificationConfig) *Factory {
	return &Factory{
		cfg: cfg,
	}
}

func (f *Factory) Build() (map[Channel]Provider, error) {
	providers := make(map[Channel]Provider)

	smsProvider, err := f.buildSMSProvider()
	if err != nil {
		return nil, err
	}
	if smsProvider != nil {
		providers[ChannelSMS] = smsProvider
	}

	emailProvider, err := f.buildEmailProvider()
	if err != nil {
		return nil, err
	}
	if emailProvider != nil {
		providers[ChannelEmail] = emailProvider
	}

	telegramProvider, err := f.buildTelegramProvider()
	if err != nil {
		return nil, err
	}
	if telegramProvider != nil {
		providers[ChannelTelegram] = telegramProvider
	}

	whatsappProvider, err := f.buildWhatsAppProvider()
	if err != nil {
		return nil, err
	}
	if whatsappProvider != nil {
		providers[ChannelWhatsApp] = whatsappProvider
	}

	return providers, nil
}

type smsAdapter struct {
	provider *sms.AfricasTalkingProvider
}

func (a *smsAdapter) Send(ctx context.Context, req Request) error {
	return a.provider.Send(ctx, sms.Request{
		Recipient: req.Recipient,
		Message:   req.Message,
	})
}

type emailAdapter struct {
	provider *email.SMTPProvider
}

func (a *emailAdapter) Send(ctx context.Context, req Request) error {
	return a.provider.Send(ctx, email.Request{
		Recipient: req.Recipient,
		Subject:   req.Subject,
		Message:   req.Message,
	})
}

type telegramAdapter struct {
	provider *telegram.TelegramProvider
}

func (a *telegramAdapter) Send(ctx context.Context, req Request) error {
	return a.provider.Send(ctx, telegram.Request{
		Recipient: req.Recipient,
		Message:   req.Message,
	})
}

type whatsappAdapter struct {
	provider *whatsapp.MetaProvider
}

func (a *whatsappAdapter) Send(ctx context.Context, req Request) error {
	return a.provider.Send(ctx, whatsapp.Request{
		Recipient: req.Recipient,
		Message:   req.Message,
	})
}

func (f *Factory) buildSMSProvider() (Provider, error) {
	if f.cfg.SMS.Provider == "" {
		return nil, nil
	}

	switch f.cfg.SMS.Provider {
	case "africas_talking":
		provider, err := sms.NewAfricasTalkingProvider(
			f.cfg.SMS,
			f.cfg.Timeout,
		)
		if err != nil {
			return nil, err
		}
		return &smsAdapter{provider: provider}, nil

	case "twilio":
		return nil, fmt.Errorf(
			"twilio SMS provider is not implemented yet",
		)

	default:
		return nil, fmt.Errorf(
			"unsupported SMS provider %q",
			f.cfg.SMS.Provider,
		)
	}
}

func (f *Factory) buildEmailProvider() (Provider, error) {
	if f.cfg.Email.Provider == "" {
		return nil, nil
	}

	switch f.cfg.Email.Provider {
	case "smtp":
		provider, err := email.NewSMTPProvider(f.cfg.Email)
		if err != nil {
			return nil, err
		}
		return &emailAdapter{provider: provider}, nil

	default:
		return nil, fmt.Errorf(
			"unsupported email provider %q",
			f.cfg.Email.Provider,
		)
	}
}

func (f *Factory) buildTelegramProvider() (Provider, error) {
	if f.cfg.Telegram.Provider == "" {
		return nil, nil
	}

	switch f.cfg.Telegram.Provider {
	case "telegram":
		provider, err := telegram.NewTelegramProvider(f.cfg.Telegram)
		if err != nil {
			return nil, err
		}
		return &telegramAdapter{provider: provider}, nil

	default:
		return nil, fmt.Errorf(
			"unsupported Telegram provider %q",
			f.cfg.Telegram.Provider,
		)
	}
}

func (f *Factory) buildWhatsAppProvider() (Provider, error) {
	if f.cfg.WhatsApp.Provider == "" {
		return nil, nil
	}

	switch f.cfg.WhatsApp.Provider {
	case "meta":
		provider, err := whatsapp.NewMetaProvider(f.cfg.WhatsApp)
		if err != nil {
			return nil, err
		}
		return &whatsappAdapter{provider: provider}, nil

	default:
		return nil, fmt.Errorf(
			"unsupported WhatsApp provider %q",
			f.cfg.WhatsApp.Provider,
		)
	}
}

var _ = time.Second
