package postgres

import (
	"chat/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type chatWMembersModel struct {
	MembersRaw []byte `db:"members"`
	chatModel
}

type memberModel struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	JoinedAt  time.Time `json:"joined_at"`
}

type chatModel struct {
	ID        uuid.UUID  `db:"chat_id"`
	Owner     *uuid.UUID `db:"chat_owner"`
	Name      string     `db:"chat_name"`
	ChatType  string     `db:"chat_type"`
	CreatedAt time.Time  `db:"created_at"`
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

func (r *ChatRepository) getSource(ctx context.Context) pgxQueryer {
	var source pgxQueryer = r.pool
	if tx, ok := ExtractTx(ctx); ok {
		source = tx
	}
	return source
}

func (r *ChatRepository) Create(ctx context.Context, owner uuid.UUID, chatName string, chatType domain.ChatType, memberIDs []uuid.UUID) (*domain.Chat, error) {

	source := r.getSource(ctx)

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
	err := source.QueryRow(ctx, query, dbObj.Owner, dbObj.Name, dbObj.ChatType).Scan(&dbObj.ID, &dbObj.CreatedAt)
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
);`, dbObj.Owner, targetUserID).Scan(&exists)
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
	_, err = source.Exec(ctx, queryAdd, dbObj.ID, dbObj.Owner)
	if err != nil {
		return nil, err
	}

	return dbObj.toDomain(), nil
}
func (r *ChatRepository) Delete(ctx context.Context, chatID uuid.UUID) error {
	source := r.getSource(ctx)
	const query = "DELETE FROM chats WHERE chat_id = $1"

	t, err := source.Exec(ctx, query, chatID)
	if err != nil {
		return err
	}

	if t.RowsAffected() != 1 {
		return domain.ErrNoAffectedRows
	}
	return nil
}

// TODO:
// Send to broker info about rename, to rename UI on every currently connected user's chat
func (r *ChatRepository) Rename(ctx context.Context, chatID uuid.UUID, newName string) error {
	source := r.getSource(ctx)
	const query = "UPDATE chats SET chat_name = $1 WHERE chat_id = $2"

	t, err := source.Exec(ctx, query, newName, chatID)
	if err != nil {
		return err
	}

	if t.RowsAffected() != 1 {
		return domain.ErrNoAffectedRows
	}
	return nil
}

func (r *ChatRepository) GetChatsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error) {
	source := r.getSource(ctx)
	const query = `
		SELECT c.chat_id, c.chat_owner, c.chat_name, c.created_at, c.chat_type 
		FROM chat_members cm 
		JOIN chats c ON cm.chat_id = c.chat_id 
		WHERE cm.user_id = $1`

	rows, err := source.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tempChats, err := pgx.CollectRows(rows, pgx.RowToStructByName[chatModel])
	if err != nil {
		return nil, err
	}

	chats := make([]*domain.Chat, 0, len(tempChats))
	for _, c := range tempChats {
		chats = append(chats, c.toDomain())
	}
	return chats, nil
}
func (r *ChatRepository) GetChatsWithMembersByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ChatWithMembers, error) {
	source := r.getSource(ctx)
	const query = `
		SELECT 
			c.chat_id, 
			c.chat_name, 
			c.chat_type, 
			c.chat_owner,
			c.created_at,
			COALESCE(
				jsonb_agg(
					jsonb_build_object(
						'user_id', u.user_id,
						'username', u.username,
						'created_at', u.created_at,
						'joined_at', cm.joined_at
					)
				), '[]'::jsonb
			) AS members
		FROM chats c
		JOIN chat_members cm ON c.chat_id = cm.chat_id
		JOIN users u ON cm.user_id = u.user_id
		WHERE c.chat_id IN (SELECT chat_id FROM chat_members WHERE user_id = $1)
		GROUP BY c.chat_id, c.chat_name, c.chat_type, c.chat_owner, c.created_at
		ORDER BY c.created_at DESC;`

	rows, err := source.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tempChats, err := pgx.CollectRows(rows, pgx.RowToStructByName[chatWMembersModel])
	if err != nil {
		return nil, fmt.Errorf("collecting chats for user %s error: %v", userID.String(), err)
	}

	chats := make([]*domain.ChatWithMembers, 0, len(tempChats))
	for _, c := range tempChats {
		cwm := domain.NewChatWithMembers(c.toDomain(), []*domain.ChatMember{})
		// Ручной парсинг JSON-массива байт в слайс структур Go
		var jsonMembers []*memberModel
		if err := json.Unmarshal(c.MembersRaw, &jsonMembers); err != nil {
			log.Printf("failed to unmarshal members json: %v", err)
		}

		for _, m := range jsonMembers {
			cwm.AddMember(domain.NewUser(m.UserID, m.Username, "", m.CreatedAt), m.JoinedAt)
		}
		chats = append(chats, cwm)
	}
	return chats, nil
}
func (r *ChatRepository) GetChatByID(ctx context.Context, chatID uuid.UUID) (*domain.Chat, error) {
	source := r.getSource(ctx)
	const query = "SELECT chat_owner, chat_name, created_at, chat_type FROM chats WHERE chat_id = $1"

	chat := chatModel{ID: chatID}
	err := source.QueryRow(ctx, query, chatID).Scan(&chat.Owner, &chat.Name, &chat.CreatedAt, &chat.ChatType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return chat.toDomain(), nil
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
func (r *ChatRepository) GetChatWithMembersByID(ctx context.Context, chatID uuid.UUID) (*domain.ChatWithMembers, error) {
	panic("not implemented")
}
