package reservation

import "errors"

var (
	ErrNoTablesAreAvailable = errors.New("no tables are available")
	ErrReservationNotFound  = errors.New("reservation not found")
)
