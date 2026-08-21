package notification

import "context"

type Provider interface {
	Send(
		ctx context.Context,
		req Request,
	) error
}
