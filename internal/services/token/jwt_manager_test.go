package token_test

import (
	"errors"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/services/token"
	"github.com/stretchr/testify/require"
)

func TestJWTManager(t *testing.T) {
	jwtManager, err := token.NewJWTManger(faker.Sentence())
	require.NoError(t, err)
	require.NotEmpty(t, jwtManager)

	userID := 1
	duration := time.Minute
	tokenString, err := jwtManager.GenerateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := jwtManager.VerifyToken(tokenString)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, userID, payload.UserID)
	require.WithinDuration(t, time.Now(), payload.IssuedAt.Time, time.Second)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiresAt.Time, time.Second)
}

func TestJWTManagerWithExpiredToken(t *testing.T) {
	jwtManager, err := token.NewJWTManger(faker.Sentence())
	require.NoError(t, err)
	require.NotEmpty(t, jwtManager)

	tokenString, err := jwtManager.GenerateToken(1, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := jwtManager.VerifyToken(tokenString)
	require.Error(t, err)
	require.True(t, errors.Is(err, jwt.ErrTokenExpired))
	require.Empty(t, payload)
}

func TestJWTManagerInvalidToken(t *testing.T) {
	jwtManager, err := token.NewJWTManger(faker.Sentence())
	require.NoError(t, err)
	require.NotEmpty(t, jwtManager)

	payload, err := token.NewPayload(1, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	tokenString, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = jwtManager.VerifyToken(tokenString)
	require.Error(t, err)
	require.Empty(t, payload)
	require.True(t, errors.Is(err, token.ErrInvalidToken))
}
