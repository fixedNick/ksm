package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChatType int

const (
	CHAT_TYPE_USER_TO_USER    ChatType = 1
	CHAT_TYPE_PUBLIC          ChatType = 2
	CHAT_TYPE_PRIVATE         ChatType = 3
	CHAT_TYPE_PUBLIC_CHANNEL  ChatType = 4
	CHAT_TYPE_PRIVATE_CHANNEL ChatType = 5
)

var mappedChatTypes = map[ChatType]string{
	CHAT_TYPE_USER_TO_USER:    "user",
	CHAT_TYPE_PUBLIC:          "public",
	CHAT_TYPE_PRIVATE:         "private",
	CHAT_TYPE_PUBLIC_CHANNEL:  "public_channel",
	CHAT_TYPE_PRIVATE_CHANNEL: "private_channel",
}

func (ct ChatType) ToSQLValue() string { return mappedChatTypes[ct] }

type Chat struct {
	id        uuid.UUID
	owner     *uuid.UUID
	name      string
	createdAt time.Time
	chatType  ChatType
}

func (c *Chat) Id() uuid.UUID        { return c.id }
func (c *Chat) Owner() *uuid.UUID    { return c.owner }
func (c *Chat) CreatedAt() time.Time { return c.createdAt }
func (c *Chat) Name() string         { return c.name }

func NewChat(id, owner uuid.UUID, name string, createdAt time.Time, chatType ChatType) *Chat {
	return &Chat{
		id:        id,
		owner:     &owner,
		name:      name,
		createdAt: createdAt,
		chatType:  chatType,
	}
}
