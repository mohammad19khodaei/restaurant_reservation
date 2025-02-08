package actions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/reservation"
)

// CancelRequest represents the request body for canceling
type CancelRequest struct {
	ID int `json:"id" binding:"required"`
}

// CancelAction is a function that handles the cancel action
func CancelAction(repository reservation.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody CancelRequest
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repository.CancelReservation(ctx, requestBody.ID); err != nil {
			if errors.Is(err, reservation.ErrReservationNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Reservation canceled successfully"})
	}
}
