package message

import (
	"context"
	"ksm-chat/internal/domain"

	"github.com/google/uuid"
)

type MessageService interface {
	SendMessage(ctx context.Context, chatID, senderID uuid.UUID, content string, attachments []string) (*domain.Message, error)
	ReadMessage(ctx context.Context, messageID uint64, callerUserID uuid.UUID) error
	EditMessage(ctx context.Context, messageID uint64, callerUserID uuid.UUID, newContent string, newAttachments []string) (*domain.EditedMessage, error)
	DeleteMessage(ctx context.Context, messageID uint64, callerUserID uuid.UUID) error
	GetChatHistory(ctx context.Context, chatID, callerUserID uuid.UUID, limit int, beforeMessageID uint64) ([]*domain.Message, error)
}
