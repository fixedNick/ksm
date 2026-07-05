package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	id          uint64
	chatID      uuid.UUID
	senderID    uuid.UUID
	content     string
	attachments []string
	createdAt   time.Time
}

type EditedMessage struct {
	*Message
	updatedAt time.Time
}

func NewMessage(id uint64, chatID, senderID uuid.UUID, content string, attachments []string, createdAt time.Time) *Message {
	return &Message{
		id:          id,
		chatID:      chatID,
		senderID:    senderID,
		content:     content,
		attachments: attachments,
		createdAt:   createdAt,
	}
}

func NewEditedMessage(id uint64, chatID, senderID uuid.UUID, content string, attachments []string, createdAt, updatedAt time.Time) *EditedMessage {
	return &EditedMessage{
		Message:   NewMessage(id, chatID, senderID, content, attachments, createdAt),
		updatedAt: updatedAt,
	}
}

func (em *EditedMessage) UpdatedAt() time.Time {
	return em.updatedAt
}

func (m *Message) ID() uint64            { return m.id }
func (m *Message) ChatID() uuid.UUID     { return m.chatID }
func (m *Message) SenderID() uuid.UUID   { return m.senderID }
func (m *Message) Content() string       { return m.content }
func (m *Message) Attachments() []string { return m.attachments }
func (m *Message) CreatedAt() time.Time  { return m.createdAt }
