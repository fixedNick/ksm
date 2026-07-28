package service

import (
	"chat/internal/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ChatService struct {
	chatRepo IChatRepository
}

func NewChatService(chatRepo IChatRepository) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

func (cs *ChatService) GetUserChats(ctx context.Context, requestFromUserID, targetID uuid.UUID) ([]*domain.Chat, error) {
	if requestFromUserID != targetID {
		return nil, fmt.Errorf("access denied")
	}
	return cs.chatRepo.GetChatsByUserID(ctx, targetID)
}
func (cs *ChatService) GetUserChatsWithMembers(ctx context.Context, requestFromUserID, targetID uuid.UUID) ([]*domain.ChatWithMembers, error) {
	if requestFromUserID != targetID {
		return nil, fmt.Errorf("access denied")
	}

	return cs.chatRepo.GetChatsWithMembersByUserID(ctx, targetID)
}

// Create new chat user-user
// Repository checks dependencies of user1->user2/user2->user1 id db, if exist - return error
func (cs *ChatService) CreatePersonalChat(ctx context.Context, ownerID, targetUserID uuid.UUID) (*domain.Chat, error) {
	return cs.chatRepo.Create(ctx, ownerID, "", domain.CHAT_TYPE_USER_TO_USER)
}

// Create new group chat, except type CHAT_TYPER_USER_TO_USER
func (cs *ChatService) CreateGroupChat(ctx context.Context, ownerID uuid.UUID, chatName string, chatType domain.ChatType, membersToAdd []uuid.UUID) (*domain.Chat, error) {
	if chatType == domain.CHAT_TYPE_USER_TO_USER {
		return nil, fmt.Errorf("cant create group chat with [user-to-user] type")
	}

	if chatName == "" {
		return nil, fmt.Errorf("group chat name cannot be empty")
	}
	return cs.chatRepo.Create(ctx, ownerID, chatName, chatType)
}

func (cs *ChatService) AddMember(ctx context.Context, chatID uuid.UUID, requestFromUserID, targetUserID uuid.UUID) error {
	return cs.addRemoveUserFromChat(ctx, true, chatID, requestFromUserID, targetUserID)
}
func (cs *ChatService) RemoveMember(ctx context.Context, chatID uuid.UUID, requestFromUserID, targetUserID uuid.UUID) error {
	return cs.addRemoveUserFromChat(ctx, false, chatID, requestFromUserID, targetUserID)
}

func (cs *ChatService) addRemoveUserFromChat(ctx context.Context, add bool, chatID, requestFromUserID, targetUserID uuid.UUID) error {
	chat, err := cs.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.Owner() == nil {
		return fmt.Errorf("access denied, chat don't have owner")
	}

	if *chat.Owner() != requestFromUserID {
		return fmt.Errorf("access denied")
	}

	inChat, err := cs.chatRepo.IsUserInChat(ctx, chatID, targetUserID)
	if err != nil {
		return err
	}

	if add && inChat {
		return fmt.Errorf("user already in chat")
	} else if !add && !inChat {
		return fmt.Errorf("user not in chat")
	}

	if add {
		return cs.chatRepo.AddMember(ctx, chatID, targetUserID)
	}

	return cs.chatRepo.RemoveMember(ctx, chatID, targetUserID)
}

func (cs *ChatService) ChangeChatName(ctx context.Context, chatID, requestFromUserID uuid.UUID, newName string) error {
	chat, err := cs.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.Owner() == nil {
		return fmt.Errorf("access denied, chat don't have owner")
	}

	if *chat.Owner() != requestFromUserID {
		return fmt.Errorf("access denied")
	}

	return cs.chatRepo.Rename(ctx, chatID, newName)
}
