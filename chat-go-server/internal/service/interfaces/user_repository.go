package interfaces

import (
	"context"
	"ksm-chat/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, username, pwdHash string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByCredentials(ctx context.Context, username, pwdHash string) (*domain.User, error)
}
