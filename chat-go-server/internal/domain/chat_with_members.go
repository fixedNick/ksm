package domain

type ChatWithMembers struct {
	chat    *Chat
	members []*ChatMember
}

func NewChatWithMembers(chat *Chat, members []*ChatMember) *ChatWithMembers {
	return &ChatWithMembers{
		chat:    chat,
		members: members,
	}
}
