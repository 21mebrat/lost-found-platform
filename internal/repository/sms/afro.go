package sms

type AfroMessageRepository struct {
	Token    string
	SenderID string
	Url      string
}

func NewAfroMessageRepository(token, senderID, url string) *AfroMessageRepository {
	return &AfroMessageRepository{
		Token:    token,
		SenderID: senderID,
		Url:      url,
	}
}
