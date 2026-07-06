package service

import (
	"context"
	"ksm-chat/internal/domain"
	"time"

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

type MessageRepository interface {
	Save(ctx context.Context, chatID, senderID uuid.UUID, content string, attachments []string, createdAt time.Time) (*domain.Message, error)
	EditMessage(ctx context.Context, messageID uint64, content string, attachments []string) (*domain.EditedMessage, error)
	GetHistory(ctx context.Context, chatID uuid.UUID, limit int, beforeMessageID uint64) ([]*domain.Message, error)

	// Mark message as read for user
	MarkAsRead(ctx context.Context, messageID uint64, userID uuid.UUID) error
	// Get user id's who reads a message
	GetReadStatus(ctx context.Context, messageID uint64) ([]uuid.UUID, error)
}

type UserRepository interface {
	Create(ctx context.Context, username, pwdHash string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type PresenceRepository interface {
	SetOnline(ctx context.Context, userID uuid.UUID, ttl time.Duration) error
	SetOffline(ctx context.Context, userID uuid.UUID) error
	GetOnlineStatus(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}

type MessageBroker interface {
	Publish(ctx context.Context, channel string, msg *domain.Message) error
	Subscribe(ctx context.Context, channel string) (<-chan *domain.Message, error)
}

type IAuthService interface {
	SignUp(ctx context.Context, username, plainPassword string) (*domain.User, error)
	SignIn(ctx context.Context, username, plainPassword string) (*domain.User, *domain.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error)
}

type IChatService interface {
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error)
	GetUserChatsWithMembers(ctx context.Context, userID uuid.UUID) ([]*domain.ChatWithMembers, error)

	CreatePersonalChat(ctx context.Context, ownerID, targetUserID uuid.UUID) (*domain.Chat, error)
	CreateGroupChat(ctx context.Context, ownerID uuid.UUID, chatType domain.ChatType, membersToAdd []uuid.UUID) (*domain.Chat, error)

	AddMember(ctx context.Context, chatID uuid.UUID, targetUserID uuid.UUID) error
	RemoveMember(ctx context.Context, chatID uuid.UUID, targetUserID uuid.UUID) error

	ChangeChatName(ctx context.Context, chatID uuid.UUID, newName string) (*domain.Chat, error)
}

type IMessageService interface {
	SendMessage(ctx context.Context, chatID, senderID uuid.UUID, content string, attachments []string) (*domain.Message, error)
	ReadMessage(ctx context.Context, messageID uint64, callerUserID uuid.UUID) error
	EditMessage(ctx context.Context, messageID uint64, callerUserID uuid.UUID, newContent string, newAttachments []string) (*domain.EditedMessage, error)
	DeleteMessage(ctx context.Context, messageID uint64, callerUserID uuid.UUID) error
	GetChatHistory(ctx context.Context, chatID, callerUserID uuid.UUID, limit int, beforeMessageID uint64) ([]*domain.Message, error)
}
