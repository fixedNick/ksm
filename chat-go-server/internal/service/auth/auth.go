package auth

import (
	"context"
	"ksm-chat/internal/domain"
)

type AuthService interface {
	SignUp(ctx context.Context, username, pwdHash string) (*domain.User, error)
	SignIn(ctx context.Context, username, pwdHash string) (*domain.User, *domain.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error)
}
