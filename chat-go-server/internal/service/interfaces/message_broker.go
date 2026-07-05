package interfaces

import (
	"context"
	"ksm-chat/internal/domain"
)

type MessageBroker interface {
	Publish(ctx context.Context, channel string, msg *domain.Message) error
	Subscribe(ctx context.Context, channel string) (<-chan *domain.Message, error)
}
