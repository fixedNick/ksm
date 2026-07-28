package postgres

import (
	"chat/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type chatModel struct {
	ID        uuid.UUID
	Owner     *uuid.UUID
	Name      string
	ChatType  string
	CreatedAt time.Time
}

func (cm *chatModel) toDomain() *domain.Chat {
	return domain.NewChat(
		cm.ID,
		cm.Owner,
		cm.Name,
		cm.CreatedAt,
		domain.ParseChatType(cm.ChatType),
	)
}

type ChatRepository struct {
	pool *pgxpool.Pool
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository { return &ChatRepository{pool: pool} }

func (r *ChatRepository) Create(ctx context.Context, owner uuid.UUID, chatName string, chatType domain.ChatType, memberIDs []uuid.UUID) (*domain.Chat, error) {
	var source pgxQueryer = r.pool
	if tx, ok := ExtractTx(ctx); ok {
		source = tx
	}

	dbObj := chatModel{
		Owner:    &owner,
		Name:     chatName,
		ChatType: chatType.ToSQLValue(),
	}

	if chatType == domain.CHAT_TYPE_USER_TO_USER {
		var targetUserID uuid.UUID = uuid.Nil
		for _, id := range memberIDs {
			if id != owner {
				targetUserID = id
				break
			}
		}
		if targetUserID == uuid.Nil {
			return nil, domain.ErrPersonalChatWithoutTargetUser
		}

		return r.createPersonalChat(ctx, source, dbObj, targetUserID)
	}

	// create group chat
	const query = "INSERT INTO chats (chat_owner, chat_name, chat_type) VALUES ($1, $2, $3) RETURNING chat_id, created_at"
	err := source.QueryRow(ctx, query, *dbObj.Owner, dbObj.Name, dbObj.ChatType).Scan(&dbObj.ID, &dbObj.CreatedAt)
	if err != nil {
		return nil, err
	}

	for _, id := range memberIDs {
		const query = "INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2)"
		_, err = source.Exec(ctx, query, dbObj.ID, id)
		if err != nil {
			return nil, err
		}
	}

	return dbObj.toDomain(), nil
}

func (r *ChatRepository) createPersonalChat(ctx context.Context, source pgxQueryer, dbObj chatModel, targetUserID uuid.UUID) (*domain.Chat, error) {
	var exists bool
	err := source.QueryRow(ctx,
		`SELECT EXISTS (
SELECT 1 
FROM chat_members cm1
JOIN chat_members cm2 ON cm1.chat_id = cm2.chat_id
JOIN chats c ON cm1.chat_id = c.chat_id
WHERE cm1.user_id = $1  
AND cm2.user_id = $2  
AND c.chat_type = 'user'::chat_type_id
);`,
	).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrAlreadyExists
	}

	const query = "INSERT INTO chats (chat_owner, chat_name, chat_type) VALUES ($1, $2, $3) RETURNING chat_id, created_at"
	err = source.QueryRow(ctx, query, dbObj.Owner, dbObj.Name, dbObj.ChatType).Scan(&dbObj.ID, &dbObj.CreatedAt)
	if err != nil {
		return nil, err
	}

	const queryAdd = "INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2)"
	_, err = source.Exec(ctx, queryAdd, dbObj.ID, targetUserID)
	if err != nil {
		return nil, err
	}
	_, err = source.Exec(ctx, queryAdd, dbObj.ID, *dbObj.Owner)
	if err != nil {
		return nil, err
	}

	return dbObj.toDomain(), nil
}
func (r *ChatRepository) Delete(ctx context.Context, chatID uuid.UUID) error {
	panic("not implemented")
}
func (r *ChatRepository) Rename(ctx context.Context, chatID uuid.UUID, newName string) error {
	panic("not implemented")
}
func (r *ChatRepository) GetChatsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error) {
	panic("not implemented")
}
func (r *ChatRepository) GetChatsWithMembersByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ChatWithMembers, error) {
	panic("not implemented")
}
func (r *ChatRepository) GetChatByID(ctx context.Context, chatID uuid.UUID) (*domain.Chat, error) {
	panic("not implemented")
}
func (r *ChatRepository) GetChatWithMembersByID(ctx context.Context, chatID uuid.UUID) (*domain.ChatWithMembers, error) {
	panic("not implemented")
}
func (r *ChatRepository) AddMember(ctx context.Context, chatID uuid.UUID, targetUserID uuid.UUID) error {
	panic("not implemented")
}
func (r *ChatRepository) RemoveMember(ctx context.Context, chatID uuid.UUID, targetUserID uuid.UUID) error {
	panic("not implemented")
}
func (r *ChatRepository) IsUserInChat(ctx context.Context, chatID uuid.UUID, targetUserID uuid.UUID) (bool, error) {
	panic("not implemented")
}
