package notification

import (
	"context"
	"errors"
	"fmt"
)

type Service struct {
	providers map[Channel]Provider
}

var _ Provider = (*Service)(nil)

func NewService(providers map[Channel]Provider) *Service {
	return &Service{
		providers: providers,
	}
}

func (s *Service) Send(
	ctx context.Context,
	req Request,
) error {
	if err := validateRequest(req); err != nil {
		return err
	}

	provider, ok := s.providers[req.Channel]
	if !ok {
		return fmt.Errorf(
			"notification provider not configured for channel %q",
			req.Channel,
		)
	}

	if err := provider.Send(ctx, req); err != nil {
		return fmt.Errorf(
			"send %s notification: %w",
			req.Channel,
			err,
		)
	}

	return nil
}

func validateRequest(req Request) error {
	if req.Channel == "" {
		return errors.New("notification channel is required")
	}

	if req.Recipient == "" {
		return errors.New("notification recipient is required")
	}

	if req.Message == "" {
		return errors.New("notification message is required")
	}

	return nil
}
