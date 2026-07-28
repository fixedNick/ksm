package service

import (
	"chat/internal/domain"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type MessageService struct {
	broker      IMessageBroker
	messageRepo IMessageRepository
	chatRepo    IChatRepository
}

func NewMessageService(broker IMessageBroker, messageRepo IMessageRepository, chatRepo IChatRepository) *MessageService {
	return &MessageService{
		broker:      broker,
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
	}
}

func (ms *MessageService) SendMessage(ctx context.Context, chatID, senderID uuid.UUID, content string, attachments []string) (*domain.Message, error) {
	// Is user in chat
	// If chat is CHANNEL type, check that user is owner
	ok, err := ms.chatRepo.IsUserInChat(ctx, chatID, senderID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("access denied: user not in chat")
	}

	chat, err := ms.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if chat.ChatType() == domain.CHAT_TYPE_PRIVATE_CHANNEL || chat.ChatType() == domain.CHAT_TYPE_PUBLIC_CHANNEL {
		if chat.Owner() == nil {
			return nil, fmt.Errorf("access denied: chat don't have owner")
		}

		if *chat.Owner() != senderID {
			return nil, fmt.Errorf("access denied: sender is not chat owner")
		}
	}

	msg, err := ms.messageRepo.Save(ctx, chatID, senderID, content, attachments, time.Now())
	if err != nil {
		return nil, err
	}

	if err = ms.broker.Publish(ctx, fmt.Sprintf("chat:%s", chatID.String()), msg); err != nil {
		// TODO: Logging
		log.Print("Error on Publish into Broker on SendMessage: " + err.Error())
	}

	return msg, nil

}
func (ms *MessageService) ReadMessage(ctx context.Context, messageID uint64, requestFromUserID uuid.UUID) error {
	if _, err := ms.isUserInChat(ctx, messageID, requestFromUserID); err != nil {
		return err
	}
	return ms.messageRepo.MarkAsRead(ctx, messageID, requestFromUserID)
}
func (ms *MessageService) EditMessage(ctx context.Context, messageID uint64, requestFromUserID uuid.UUID, newContent string, newAttachments []string) (*domain.EditedMessage, error) {
	msg, err := ms.isUserInChat(ctx, messageID, requestFromUserID)
	if err != nil {
		return nil, err
	}

	// is msg owner == request user
	if msg.SenderID() != requestFromUserID {
		return nil, fmt.Errorf("access denied: only owner possible to edit message")
	}

	return ms.messageRepo.Edit(ctx, messageID, newContent, newAttachments)
}
func (ms *MessageService) DeleteMessage(ctx context.Context, messageID uint64, requestFromUserID uuid.UUID) error {
	msg, err := ms.isUserInChat(ctx, messageID, requestFromUserID)
	if err != nil {
		return err
	}

	// is msg owner == request user
	if msg.SenderID() != requestFromUserID {
		return fmt.Errorf("access denied: only owner possible to delete message")
	}

	_, err = ms.messageRepo.Edit(ctx, messageID, "Message deleted", []string{})
	return err
}
func (ms *MessageService) GetChatHistory(ctx context.Context, chatID, requestFromUserID uuid.UUID, limit int, beforeMessageID uint64) ([]*domain.Message, error) {
	ok, err := ms.chatRepo.IsUserInChat(ctx, chatID, requestFromUserID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("access denied: user not in chat")
	}

	return ms.messageRepo.GetHistory(ctx, chatID, limit, beforeMessageID)
}

// Receive chat by messageID
// Check is User in chat
// return message and error
func (ms *MessageService) isUserInChat(ctx context.Context, messageID uint64, requestFromUserID uuid.UUID) (*domain.Message, error) {
	msg, err := ms.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	ok, err := ms.chatRepo.IsUserInChat(ctx, msg.ChatID(), requestFromUserID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("access denied: user not in chat")
	}
	return msg, nil
}
