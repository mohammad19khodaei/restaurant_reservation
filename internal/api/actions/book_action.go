package actions

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/api/middlewares"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/reservation"
)

// BookRequest represents the request body for booking
type BookRequest struct {
	SeatsCount int    `json:"seats_count" binding:"required,min=1,max=10"`
	Date       string `json:"date" binding:"required"`
}

// BookResponse represents the response body for booking
type BookResponse struct {
	ID         int     `json:"id"`
	TableID    int     `json:"table_id"`
	SeatsCount int     `json:"seats_count"`
	Price      float64 `json:"price"`
}

// BookAction is a function that handles the book action
func BookAction(reservationRepo reservation.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody BookRequest
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate date format (YYYY-MM-DD)
		date, err := time.Parse("2006-01-02", requestBody.Date)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected YYYY-MM-DD"})
			return
		}

		if time.Now().After(date) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date, should be in the future"})
			return
		}

		seatsCount := requestBody.SeatsCount
		if seatsCount%2 != 0 {
			seatsCount++
		}

		userID := ctx.MustGet(middlewares.AuthUserIDKey).(int)

		resv, err := reservationRepo.BookTable(ctx, userID, seatsCount, date)
		if err != nil {
			if errors.Is(err, reservation.ErrNoTablesAreAvailable) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res := BookResponse{
			ID:         int(resv.ID),
			TableID:    int(resv.TableID),
			SeatsCount: resv.SeatsCount,
			Price:      resv.Price,
		}
		ctx.JSON(http.StatusOK, res)
	}
}
