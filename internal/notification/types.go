package notification

// Channel represents the communication channel used
// to deliver a notification.
type Channel string

const (
	ChannelSMS      Channel = "sms"
	ChannelEmail    Channel = "email"
	ChannelTelegram Channel = "telegram"
	ChannelWhatsApp Channel = "whatsapp"
)

// Request contains everything required to send
// a notification through a provider.
type Request struct {
	Channel   Channel
	Recipient string
	Subject   string
	Message   string
}
