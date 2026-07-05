package interfaces

import (
	"context"
	"ksm-chat/internal/domain"

	"github.com/google/uuid"
)

type ChatRepository interface {
	Create(ctx context.Context, owner uuid.UUID, chatName string, chatType domain.ChatType) (*domain.Chat, error)
	Delete(ctx context.Context, chatID uuid.UUID) error
	Rename(ctx context.Context, chatID uuid.UUID, newName string) error
	GetChatsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error)
	GetChatsWithMembersByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ChatWithMembers, error)
	GetChatByID(ctx context.Context, chatID uuid.UUID) (*domain.Chat, error)
}
