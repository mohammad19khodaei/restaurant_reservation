package token

import "time"

type Manager interface {
	GenerateToken(userID int, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
