package jwt

import (
	"crypto/rand"
	"ksm-chat/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secret         []byte
	accessTokenTTL time.Duration
}

type AccessClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTService(sectetLength int, accessTTL time.Duration) *JWTService {
	return &JWTService{
		secret:         generateSecret(sectetLength),
		accessTokenTTL: accessTTL,
	}
}

func generateSecret(len int) []byte {
	bytes := make([]byte, len)
	if _, err := rand.Read(bytes); err != nil {
		panic("Error on generate secret for JWT: " + err.Error())
	}
	return bytes
}

func (js *JWTService) Generate(username string, user_id uuid.UUID, issuer Issuer) (*domain.Token, error) {
	claims := AccessClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user_id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    string(issuer),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(js.secret)
	if err != nil {
		return nil, err
	}

	return domain.NewToken(ss, uuid.NewString()), nil
}
