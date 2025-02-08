package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload is a struct that holds the claims for JWT
type Payload struct {
	UserID int `json:"username"`
	jwt.RegisteredClaims
}

// NewPayload creates a new Payload
func NewPayload(userID int, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return payload, nil
}
