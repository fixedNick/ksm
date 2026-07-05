package domain

type ChatWithMembes struct {
	chat    *Chat
	members []*ChatMember
}

func NewChatWithMembers(chat *Chat, members []*ChatMember) *ChatWithMembes {
	return &ChatWithMembes{
		chat:    chat,
		members: members,
	}
}
