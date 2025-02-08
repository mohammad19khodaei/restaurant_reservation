package reservation

import (
	"context"
	"time"
)

type Repository interface {
	BookTable(ctx context.Context, userID int, seatsNeeded int, date time.Time) (*Reservation, error)
	CancelReservation(ctx context.Context, reservationID int) error
}
