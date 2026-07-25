package service

import (
	"context"
	"fmt"
	"ksm-chat/internal/domain"
	"ksm-chat/internal/infrastructure/jwt"
	"ksm-chat/internal/infrastructure/limits"
	"strings"
)

type AuthService struct {
	userRepository IUserRepository
	hasher         IHasherComparer
	jwtService     ITokenGenerator
}

func NewAuthService(userRepo IUserRepository, hasher IHasherComparer, jwtService ITokenGenerator) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		hasher:         hasher,
		jwtService:     jwtService,
	}
}

func (as *AuthService) SignUp(ctx context.Context, username, plainPassword string) (*domain.User, *domain.Token, error) {
	// validate is username and passwords valid
	if err := validateCredentials(username, plainPassword); err != nil {
		return nil, nil, err
	}

	// hash password
	pwdHash, err := as.hasher.Hash(plainPassword)
	if err != nil {
		return nil, nil, err
	}

	// Create User Account
	user, err := as.userRepository.Create(ctx, username, pwdHash)
	if err != nil {
		return nil, nil, err
	}

	// Generate JWT
	token, err := as.jwtService.Generate(user.Username(), user.ID(), string(jwt.ISSUER_AUTH))
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}
func (as *AuthService) SignIn(ctx context.Context, username, plainPassword string) (*domain.User, *domain.Token, error) {
	// Check is input valid
	if err := validateCredentials(username, plainPassword); err != nil {
		return nil, nil, err
	}

	user, err := as.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	// Compare passwords
	if err := as.hasher.Compare(plainPassword, user.PasswordHash()); err != nil {
		return nil, nil, err
	}

	// Passwords matched
	// Generate JWT token
	token, err := as.jwtService.Generate(username, user.ID(), string(jwt.ISSUER_AUTH))
	return user, token, err
}
func (as *AuthService) UpdateTokens(ctx context.Context, refreshToken string) (*domain.Token, error) {
	panic("not implemented")
}

func validateCredentials(username string, plainPassword string) error {
	if len(username) > limits.USERNAME_MAX_LENGTH || len(username) < limits.USERNAME_MIN_LENGTH {
		return fmt.Errorf("username must be in range [> %d and < %d]", limits.USERNAME_MIN_LENGTH, limits.USERNAME_MAX_LENGTH)
	}

	if len(plainPassword) > limits.PASSWORD_MAX_LENGTH || len(plainPassword) < limits.PASSWORD_MIN_LENGTH {
		return fmt.Errorf("password must be in range [> %d and < %d]", limits.PASSWORD_MIN_LENGTH, limits.PASSWORD_MAX_LENGTH)
	}

	for _, symbolsRange := range limits.PASSWORD_REQUIRE_SYMBOLS {
		symbolFound := false
		for _, symbol := range symbolsRange {
			if strings.ContainsRune(plainPassword, symbol) {
				symbolFound = true
				break
			}
		}

		if !symbolFound {
			return fmt.Errorf("password must contain one of the symbols: [%s]", symbolsRange)
		}
	}

	return nil
}
