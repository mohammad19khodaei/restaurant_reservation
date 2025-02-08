package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/reservation"
	"gorm.io/gorm"
)

// GormReservationRepository is a repository for reservation operations
type GormReservationRepository struct {
	db *gorm.DB
}

// NewGormReservationRepository creates a new instance of GormReservationRepository
func NewGormReservationRepository(db *gorm.DB) reservation.Repository {
	return &GormReservationRepository{db: db}
}

// BookTable books a table for a user on a specific date
func (r *GormReservationRepository) BookTable(ctx context.Context, userID int, seatsNeeded int, date time.Time) (*reservation.Reservation, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	var tableID uint
	var seatPrice, totalPrice float64

	query := `
		WITH table_availability AS (
			SELECT t.id AS table_id, t.seats_count AS total_seats,
				COALESCE(SUM(r.seats_count), 0) AS reserved_seats,
				(t.seats_count - COALESCE(SUM(r.seats_count), 0)) AS available_seats
			FROM tables t
			LEFT JOIN reservations r 
				ON t.id = r.table_id AND r.date = ?
			GROUP BY t.id, t.seats_count
		),
		selected_table AS (
			SELECT table_id FROM table_availability
			WHERE available_seats >= ? 
			ORDER BY available_seats ASC 
			LIMIT 1
		),
		seat_price AS (
			SELECT seat_price FROM table_settings LIMIT 1
		)
		SELECT 
			t.id AS table_id,
			sp.seat_price,
			CASE 
				WHEN ? = t.seats_count THEN (t.seats_count - 1) * sp.seat_price 
				ELSE ? * sp.seat_price
			END AS total_price
		FROM tables t
		JOIN selected_table st ON t.id = st.table_id
		JOIN seat_price sp ON true;
	`

	row := tx.Raw(query, date, seatsNeeded, seatsNeeded, seatsNeeded).Row()
	if err := row.Scan(&tableID, &seatPrice, &totalPrice); err != nil {
		tx.Rollback()
		return nil, reservation.ErrNoTablesAreAvailable
	}

	newReservation := reservation.Reservation{
		UserID:     uint(userID),
		TableID:    tableID,
		SeatsCount: seatsNeeded,
		Price:      totalPrice,
		Date:       date,
	}
	if err := tx.Create(&newReservation).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &newReservation, nil
}

// CancelReservation cancels a reservation by its ID
func (r *GormReservationRepository) CancelReservation(ctx context.Context, reservationID int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	var resv reservation.Reservation
	if err := tx.First(&resv, reservationID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reservation.ErrReservationNotFound
		}
		return err
	}

	if err := tx.Delete(&resv).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
