package domain

import "time"

type ChatWithMembers struct {
	chat    *Chat
	members []*ChatMember
}

type ChatMember struct {
	User     *User
	JoinedAt time.Time
}

func NewChatWithMembers(chat *Chat, members []*ChatMember) *ChatWithMembers {
	return &ChatWithMembers{
		chat:    chat,
		members: members,
	}
}

func (c *ChatWithMembers) Chat() *Chat {
	return c.chat
}

func (c *ChatWithMembers) Members() []*ChatMember {
	return c.members
}

func (c *ChatWithMembers) AddMember(u *User, joinedAt time.Time) {
	c.members = append(c.members, &ChatMember{
		User:     u,
		JoinedAt: joinedAt,
	})
}
