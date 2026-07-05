package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PresenceRepository interface {
	SetOnline(ctx context.Context, userID uuid.UUID, ttl time.Duration) error
	SetOffline(ctx context.Context, userID uuid.UUID) error
	GetOnlineStatus(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}
