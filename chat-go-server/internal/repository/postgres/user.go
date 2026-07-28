package postgres

import (
	"chat/internal/domain"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userModel struct {
	ID           uuid.UUID `db:"user_id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func (um *userModel) toDomain() *domain.User {
	return domain.NewUser(
		um.ID,
		um.Username,
		um.PasswordHash,
		um.CreatedAt,
	)
}

type UserRepository struct {
	dbpool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{dbpool: pool}
}

func (r *UserRepository) Create(ctx context.Context, username, passwordHash string) (*domain.User, error) {
	const query = "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING user_id, created_at"
	dbObj := userModel{
		Username:     username,
		PasswordHash: passwordHash,
	}

	var source pgxQueryer = r.dbpool
	if tx, ok := ExtractTx(ctx); ok {
		source = tx
	}

	err := source.QueryRow(ctx, query, username, passwordHash).Scan(&dbObj.ID, &dbObj.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrAlreadyExists
		}
		return nil, err
	}
	return dbObj.toDomain(), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const query = "SELECT username, password_hash, created_at FROM users WHERE user_id = $1"
	dbObj := &userModel{}

	var source pgxQueryer = r.dbpool
	if tx, ok := ExtractTx(ctx); ok {
		source = tx
	}

	err := source.QueryRow(ctx, query, id).Scan(&dbObj.Username, &dbObj.PasswordHash, &dbObj.CreatedAt)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}

	dbObj.ID = id
	return dbObj.toDomain(), err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	const query = "SELECT user_id, password_hash, created_at FROM users WHERE username = $1"
	dbObj := &userModel{}

	var source pgxQueryer = r.dbpool
	if tx, ok := ExtractTx(ctx); ok {
		source = tx
	}

	err := source.QueryRow(ctx, query, username).Scan(&dbObj.ID, &dbObj.PasswordHash, &dbObj.CreatedAt)
	dbObj.Username = username
	return dbObj.toDomain(), err
}
