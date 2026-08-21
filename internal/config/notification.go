package config

import "time"

type NotificationConfig struct {
	Timeout time.Duration

	SMS      SMSConfig
	Email    EmailConfig
	Telegram TelegramConfig
	WhatsApp WhatsAppConfig
}

type SMSConfig struct {
	Provider string

	APIKey   string
	Username string
	Password string
	SenderID string
	BaseURL  string
}

type EmailConfig struct {
	Provider string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	From         string
}

type TelegramConfig struct {
	Provider string

	BotToken string
	BaseURL  string
}

type WhatsAppConfig struct {
	Provider string

	AccessToken   string
	PhoneNumberID string
	BaseURL       string
}
