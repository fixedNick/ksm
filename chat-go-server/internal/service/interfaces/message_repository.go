package interfaces

import (
	"context"
	"ksm-chat/internal/domain"
	"time"

	"github.com/google/uuid"
)

type MessageRepository interface {
	Save(ctx context.Context, chatID, senderID uuid.UUID, content string, attachments []string, createdAt time.Time) (*domain.Message, error)
	EditMessage(ctx context.Context, messageID uint64, content string, attachments []string) (*domain.EditedMessage, error)
	GetHistory(ctx context.Context, chatID uuid.UUID, limit int, beforeMessageID uint64) ([]*domain.Message, error)

	// Mark message as read for user
	MarkAsRead(ctx context.Context, messageID uint64, userID uuid.UUID) error
	// Get user id's who reads a message
	GetReadStatus(ctx context.Context, messageID uint64) ([]uuid.UUID, error)
}
