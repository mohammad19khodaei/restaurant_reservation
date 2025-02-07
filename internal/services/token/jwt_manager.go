package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretLength = 32

// JWTManager is a struct that holds the secret key for JWT
type JWTManager struct {
	secretKey string
}

// NewJWTManger creates a new JWTManager
func NewJWTManger(secretKey string) (Manager, error) {
	if len(secretKey) < minSecretLength {
		return nil, errors.New(fmt.Sprintf("valid secret key size must be at least %d characters", minSecretLength))
	}
	return &JWTManager{
		secretKey: secretKey,
	}, nil
}

// GenerateToken generates a new JWT token
func (m *JWTManager) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(m.secretKey))
}

// VerifyToken verifies a JWT token
func (m *JWTManager) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
