package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChatMember struct {
	chatID   uuid.UUID
	userID   uuid.UUID
	joinedAt time.Time
}

func NewChatMember(chaId, userId uuid.UUID, joinedAt time.Time) *ChatMember {
	return &ChatMember{
		userID:   userId,
		chatID:   chaId,
		joinedAt: joinedAt,
	}
}

func (cm *ChatMember) ChatID() uuid.UUID   { return cm.chatID }
func (cm *ChatMember) UserID() uuid.UUID   { return cm.userID }
func (cm *ChatMember) JoinedAt() time.Time { return cm.joinedAt }
