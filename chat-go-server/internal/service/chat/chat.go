package chat

import (
	"context"
	"ksm-chat/internal/domain"

	"github.com/google/uuid"
)

type ChatService interface {
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error)
	GetUserChatsWithMembers(ctx context.Context, userID uuid.UUID) ([]*domain.ChatWithMembers, error)

	CreatePersonalChat(ctx context.Context, ownerID, targetUserID uuid.UUID) (*domain.Chat, error)
	CreateGroupChat(ctx context.Context, ownerID uuid.UUID, chatType domain.ChatType, membersToAdd []uuid.UUID) (*domain.Chat, error)

	AddMember(ctx context.Context, chatID uuid.UUID, targetUserID uuid.UUID) error
	RemoveMember(ctx context.Context, chatID uuid.UUID, targetUserID uuid.UUID) error

	ChangeChatName(ctx context.Context, chatID uuid.UUID, newName string) (*domain.Chat, error)
}
