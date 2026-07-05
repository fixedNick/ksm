package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	id           uuid.UUID
	username     string
	createdAt    time.Time
	passwordHash string
}

func (u *User) ID() uuid.UUID        { return u.id }
func (u *User) Username() string     { return u.username }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func NewUser(id uuid.UUID, username string, passwordHash string, createdAt time.Time) *User {
	return &User{id: id, username: username, createdAt: createdAt, passwordHash: passwordHash}
}
