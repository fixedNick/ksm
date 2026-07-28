package postgres

import (
	"chat/internal/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	conn pgxpool.Pool
}

// func NewUserRepository(conn *pgx.ConnPool) *UserRepository {
// 	pgx.NewConnPool()
// 	return &UserRepository{conn: conn}
// }

func (r *UserRepository) Create(ctx context.Context, username, pwdHash string) (*domain.User, error) {
	panic("not implemented")
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	panic("not implemented")
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	panic("not implemented")
}
